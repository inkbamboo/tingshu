package consts

//go:generate gonum -types=TingShuTypeEnum

type TingShuTypeEnum struct {
	NianYin string `enum:"nianYin,念音网"`
	ShuYin  string `enum:"shuYin,书音网"`
}
