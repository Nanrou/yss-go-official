package meta

type site struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	WorkTime string `json:"work_time"`
	Position string `json:"position"`
}

var Sites = [...] *site{
	{
		Name:     "香洲营业厅",
		Address:  "香洲梅华东路338号",
		Phone:    "2224838",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.284590,113.556340",
	},
	{
		Name:     "前山营业厅",
		Address:  "前山翠前南路61号",
		Phone:    "8656088",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.245719,113.522109",
	},
	{
		Name:     "金唐营业厅",
		Address:  "唐乐路18号1楼",
		Phone:    "3616605",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.360658,113.598284",
	},
	{
		Name:     "珠海市行政服务厅",
		Address:  "珠海市红山路230号",
		Phone:    "8132659",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.277020,113.541380",
	},
	{
		Name:     "拱北营业厅",
		Address:  "拱北粤海中路2083号",
		Phone:    "8885369",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.230290,113.540030",
	},
	{
		Name:     "南湾营业厅",
		Address:  "南屏坪岚路3号",
		Phone:    "8678962",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.220543,113.508324",
	},
	{
		Name:     "横琴营业厅",
		Address:  "横琴新区海河街63号华兴楼1单元104",
		Phone:    "8688091",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.111579,113.547297",
	},
	{
		Name:     "红旗营业厅",
		Address:  "金湾区红旗镇藤山一路自来水公司综合楼一楼",
		Phone:    "7791195",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.140423,113.344999",
	},
	{
		Name:     "三灶营业厅",
		Address:  "金湾区三灶镇金海岸大道东445号二楼",
		Phone:    "7761237",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.058071,113.369892",
	},
	{
		Name:     "平沙营业厅",
		Address:  "金湾区平沙镇升平大道366号1栋平沙镇党群服务中心服务厅",
		Phone:    "7751263",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.097292,113.192980",
	},
	{
		Name:     "南水营业厅",
		Address:  "南水镇南港中路128号1楼",
		Phone:    "7713225",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.025190,113.228800",
	},
	{
		Name:     "井岸营业厅",
		Address:  "井岸镇环山中路165号",
		Phone:    "5102031",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.204460,113.292810",
	},
	{
		Name:     "白蕉营业厅",
		Address:  "斗门区白藤湖湖滨一区20号",
		Phone:    "5566340",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.182080,113.330790",
	},
	{
		Name:     "桥东营业厅",
		Address:  "珠海市斗门区白蕉镇连桥路334号",
		Phone:    "5502021",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.226667,113.309921",
	},
	{
		Name:     "六乡营业厅",
		Address:  "斗门区白蕉镇六乡建设路52号",
		Phone:    "5587581",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.305893,113.284369",
	},
	{
		Name:     "斗门营业厅",
		Address:  "斗门镇斗门大道北268号",
		Phone:    "5783853",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.238740,113.197800",
	},
	{
		Name:     "乾务营业厅",
		Address:  "乾务镇盛兴3路",
		Phone:    "5581139",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.155203,113.228390",
	},
	{
		Name:     "五山营业厅",
		Address:  "乾务镇五山沙龙中路171号",
		Phone:    "5652277",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.167530,113.176860",
	},
	{
		Name:     "上横营业厅",
		Address:  "斗门区莲洲镇横山街191、193号",
		Phone:    "5562432",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.317590,113.206010",
	},
	{
		Name:     "莲洲营业厅",
		Address:  "莲洲镇莲溪圩镇河南路6号",
		Phone:    "5593863",
		WorkTime: "周一至周五 8:00-17:00",
		Position: "22.378460,113.232460",
	},
}
