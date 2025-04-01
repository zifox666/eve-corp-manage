package sde

import (
	"compress/bzip2"
	"crypto/md5"
	"encoding/hex"
	"eve-corp-manager/config"
	"eve-corp-manager/core/common"
	"eve-corp-manager/global"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SdeDownloadURL = "https://www.fuzzwork.co.uk/dump/sqlite-latest.sqlite.bz2"
	SdeMd5URL      = "https://www.fuzzwork.co.uk/dump/sqlite-latest.sqlite.bz2.md5"
	VersionFile    = "/logs/sde_version.txt"
)

// DownloadSDE 从指定URL下载SDE数据库文件
func DownloadSDE(destPath string) error {
	// 创建目标目录
	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 下载压缩文件
	compressedPath := destPath + ".bz2"
	err = downloadFile(SdeDownloadURL, compressedPath)
	if err != nil {
		return fmt.Errorf("下载SDE文件失败: %v", err)
	}

	// 解压文件
	err = decompressBzip2(compressedPath, destPath)
	if err != nil {
		return fmt.Errorf("解压SDE文件失败: %v", err)
	}

	// 更新版本信息
	err = updateVersionInfo()
	if err != nil {
		return fmt.Errorf("更新版本信息失败: %v", err)
	}

	return nil
}

// CheckMD5 检查本地SDE文件的MD5是否与远程一致
func CheckMD5(filePath string) (bool, error) {
	// 获取远程MD5
	remoteMD5, err := getRemoteMD5()
	if err != nil {
		return false, fmt.Errorf("获取远程MD5失败: %v", err)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	}

	// 计算本地文件的MD5
	localMD5, err := calculateFileMD5(filePath)
	if err != nil {
		return false, fmt.Errorf("计算本地文件MD5失败: %v", err)
	}

	// 比较MD5
	return localMD5 == remoteMD5, nil
}

// UpdateSDE 检查并更新SDE数据库
func UpdateSDE(filePath string) error {
	isValid, err := CheckMD5(filePath)
	if err != nil {
		return err
	}

	if !isValid {
		global.Logger.Info("正在更新SDE数据库...")
		err = DownloadSDE(filePath)
		if err != nil {
			return err
		}
		global.Logger.Info("SDE数据库更新完成")
	} else {
		global.Logger.Info("SDE数据库已是最新版本")
	}

	return nil
}

// 下载文件的辅助函数
func downloadFile(url, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// 解压bzip2文件
func decompressBzip2(compressedPath, destPath string) error {
	f, err := os.Open(compressedPath)
	if err != nil {
		return err
	}
	defer f.Close()

	bz2Reader := bzip2.NewReader(f)
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, bz2Reader)
	return err
}

// 获取远程MD5值
func getRemoteMD5() (string, error) {
	resp, err := http.Get(SdeMd5URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	md5str := string(body)
	// 通常MD5文件格式为: "md5sum 文件名"，提取MD5部分
	parts := strings.Fields(md5str)
	if len(parts) > 0 {
		return parts[0], nil
	}
	return "", fmt.Errorf("无效的MD5格式")
}

// 计算文件的MD5值
func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 更新版本信息文件
func updateVersionInfo() error {
	versionFilePath := config.LogsDir + VersionFile

	// 创建目录
	err := os.MkdirAll(filepath.Dir(versionFilePath), 0755)
	if err != nil {
		return err
	}

	// 获取当前时间
	now := time.Now().Format(common.TimeFormatMode1)

	// 获取远程MD5
	md5, err := getRemoteMD5()
	if err != nil {
		return err
	}

	// 写入版本信息文件
	versionInfo := fmt.Sprintf("更新时间: %s\nMD5: %s\n下载源: %s", now, md5, SdeDownloadURL)
	return os.WriteFile(versionFilePath, []byte(versionInfo), 0644)
}
