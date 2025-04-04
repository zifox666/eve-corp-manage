package utils

import (
	"eve-corp-manager/core/esi"
	"fmt"
	"strconv"
)

// PersonalLocationFlag 槽位类型枚举
const (
	HighSlot      = 1 // 高能量槽
	MediumSlot    = 2 // 中能量槽
	LowSlot       = 3 // 低能量槽
	RigSlot       = 4 // 船插
	DroneSlot     = 5 // 无人机仓
	CargoSlot     = 6 // 货柜仓
	SubsystemSlot = 7 // 子系统仓
)

// GetSlotNameByFlag 根据槽位标志获取槽位名称
func GetSlotNameByFlag(flag int) int {
	switch flag {
	case 27, 28, 29, 30, 31, 32, 33, 34: // 高槽
		return HighSlot
	case 19, 20, 21, 22, 23, 24, 25, 26: // 中槽
		return MediumSlot
	case 11, 12, 13, 14, 15, 16, 17, 18: // 低槽
		return LowSlot
	case 92, 93, 94: // 船插
		return RigSlot
	case 87: // 无人机仓
		return DroneSlot
	case 89, 5, 133, 134, 136, 137, 138: // 货柜仓
		return CargoSlot
	case 125, 126, 127, 128: // 子系统
		return SubsystemSlot
	default:
		return CargoSlot
	}
}

// KillmailItem 击毁邮件物品信息
type KillmailItem struct {
	ItemID   int    `json:"item_id"`
	ItemName string `json:"item_name"`
	ItemNum  int    `json:"item_num"`
	DropType bool   `json:"drop_type"`
	SlotType int    `json:"slot_type"`
}

// KillmailDetails 击毁邮件详情
type KillmailDetails struct {
	KillmailID      int                    `json:"killmail_id"`
	KillmailHash    string                 `json:"killmail_hash"`
	Time            string                 `json:"time"`
	SolarSystemID   int                    `json:"solar_system_id"`
	AllianceID      int                    `json:"alliance_id"`
	CorporationID   int                    `json:"corporation_id"`
	CharacterID     int                    `json:"character_id"`
	ShipTypeID      int                    `json:"ship_type_id"`
	SolarSystemName string                 `json:"solar_system_name"`
	AllianceName    string                 `json:"alliance_name"`
	CorporationName string                 `json:"corporation_name"`
	CharacterName   string                 `json:"character_name"`
	ShipTypeName    string                 `json:"ship_type_name"`
	Items           []KillmailItem         `json:"items"`
	JaniceAmount    float64                `json:"janice_amount"`
	UIJSON          map[string]interface{} `json:"ui_json"`
}

// NewKillmailDetails 创建一个新的击毁邮件详情对象
func NewKillmailDetails(killmailID int, killmailHash string) *KillmailDetails {
	return &KillmailDetails{
		KillmailID:   killmailID,
		KillmailHash: killmailHash,
		Items:        make([]KillmailItem, 0),
		UIJSON:       make(map[string]interface{}),
	}
}

// Init 初始化击毁邮件详情
func (k *KillmailDetails) Init() error {
	// 获取击毁邮件数据
	killmailData, err := esi.GetKillmail(k.KillmailID, k.KillmailHash)
	if err != nil {
		return err
	}

	// 设置基本信息
	k.Time = killmailData["killmail_time"].(string)

	victim, ok := killmailData["victim"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("无法解析victim数据")
	}

	k.SolarSystemID = int(killmailData["solar_system_id"].(float64))

	if allianceID, ok := victim["alliance_id"].(float64); ok {
		k.AllianceID = int(allianceID)
	}

	if corporationID, ok := victim["corporation_id"].(float64); ok {
		k.CorporationID = int(corporationID)
	}

	if characterID, ok := victim["character_id"].(float64); ok {
		k.CharacterID = int(characterID)
	}

	if shipTypeID, ok := victim["ship_type_id"].(float64); ok {
		k.ShipTypeID = int(shipTypeID)
	}

	// 处理物品
	items, ok := victim["items"].([]interface{})
	if ok && len(items) > 0 {
		if err := k.handleItems(items); err != nil {
			return err
		}
	}

	// 获取角色信息
	if err := k.getCharInfo(); err != nil {
		return err
	}

	// 获取Janice估价
	if err := k.getJaniceAmount(); err != nil {
		return err
	}

	// 处理UI JSON
	k.handleUIJSON()

	return nil
}

// handleItems 处理物品信息
func (k *KillmailDetails) handleItems(items []interface{}) error {
	if len(items) == 0 {
		k.Items = []KillmailItem{}
		return nil
	}

	// 提取所有物品ID
	ids := make([]int, 0, len(items))
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		if itemTypeID, ok := itemMap["item_type_id"].(float64); ok {
			ids = append(ids, int(itemTypeID))
		}
	}

	// 获取物品名称
	names, err := esi.PostIdsToNames(ids)
	if err != nil {
		return err
	}

	// 处理物品数据
	itemDict := make(map[string]*KillmailItem)
	for _, item := range items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		itemID := int(itemMap["item_type_id"].(float64))
		dropType := false
		if _, ok := itemMap["quantity_dropped"]; ok {
			dropType = true
		}

		flag := int(itemMap["flag"].(float64))
		slotType := GetSlotNameByFlag(flag)

		// 创建唯一键
		key := fmt.Sprintf("%d-%v-%d", itemID, dropType, slotType)

		if _, exists := itemDict[key]; !exists {
			itemDict[key] = &KillmailItem{
				ItemID:   itemID,
				ItemName: names[strconv.Itoa(itemID)],
				ItemNum:  0,
				DropType: dropType,
				SlotType: slotType,
			}
		}

		// 计算数量
		var quantityDropped, quantityDestroyed float64
		if qd, ok := itemMap["quantity_dropped"].(float64); ok {
			quantityDropped = qd
		}
		if qd, ok := itemMap["quantity_destroyed"].(float64); ok {
			quantityDestroyed = qd
		}

		itemDict[key].ItemNum += int(quantityDropped + quantityDestroyed)
	}

	// 转换为切片
	k.Items = make([]KillmailItem, 0, len(itemDict))
	for _, item := range itemDict {
		k.Items = append(k.Items, *item)
	}

	return nil
}

// getJaniceAmount 获取Janice估价
func (k *KillmailDetails) getJaniceAmount() error {
	queryStr := ""
	for _, item := range k.Items {
		queryStr += fmt.Sprintf("%s\t%d\n", item.ItemName, item.ItemNum)
	}
	queryStr += k.ShipTypeName

	amount, err := esi.GetAppraisal(queryStr)
	if err != nil {
		return err
	}

	k.JaniceAmount = amount
	return nil
}

// getCharInfo 获取角色信息
func (k *KillmailDetails) getCharInfo() error {
	ids := []int{
		k.CharacterID,
		k.CorporationID,
		k.AllianceID,
		k.ShipTypeID,
		k.SolarSystemID,
	}

	names, err := esi.PostIdsToNames(ids)
	if err != nil {
		return err
	}

	k.CharacterName = names[strconv.Itoa(k.CharacterID)]
	if k.CharacterName == "" {
		k.CharacterName = "未知"
	}

	k.CorporationName = names[strconv.Itoa(k.CorporationID)]
	if k.CorporationName == "" {
		k.CorporationName = "未知"
	}

	k.AllianceName = names[strconv.Itoa(k.AllianceID)]
	if k.AllianceName == "" {
		k.AllianceName = "未知"
	}

	k.ShipTypeName = names[strconv.Itoa(k.ShipTypeID)]
	if k.ShipTypeName == "" {
		k.ShipTypeName = "未知"
	}

	k.SolarSystemName = names[strconv.Itoa(k.SolarSystemID)]
	if k.SolarSystemName == "" {
		k.SolarSystemName = "未知"
	}

	return nil
}

// handleUIJSON 处理UI JSON数据
func (k *KillmailDetails) handleUIJSON() {
	k.UIJSON = map[string]interface{}{
		"killId":   k.KillmailID,
		"data":     []interface{}{},
		"value":    k.JaniceAmount,
		"charName": k.CharacterName,
		"shipName": k.ShipTypeName,
		"corpName": k.CorporationName,
		"time":     k.Time,
	}

	// 槽位映射
	slotMap := map[int]string{
		1: "高能量槽",
		2: "中能量槽",
		3: "低能量槽",
		4: "船插",
		5: "无人机仓",
		6: "货柜仓",
		7: "子系统仓",
	}

	// 按槽位分组物品
	slotData := make(map[string][]map[string]interface{})
	for _, slotName := range slotMap {
		slotData[slotName] = []map[string]interface{}{}
	}

	for _, item := range k.Items {
		slotName := slotMap[item.SlotType]
		if slotName == "" {
			slotName = "货柜仓"
		}

		slotData[slotName] = append(slotData[slotName], map[string]interface{}{
			"id":      item.ItemID,
			"name":    item.ItemName,
			"dropped": fmt.Sprintf("%v", item.DropType),
			"num":     item.ItemNum,
		})
	}

	// 按顺序添加槽位数据
	var data []map[string]interface{}
	for i := 1; i <= 7; i++ {
		slotName := slotMap[i]
		if len(slotData[slotName]) > 0 {
			data = append(data, map[string]interface{}{
				"slotName": slotName,
				"data":     slotData[slotName],
			})
		}
	}

	k.UIJSON["data"] = data
}
