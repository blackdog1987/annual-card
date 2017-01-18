package data

// 计算分页数量
func CalcPage(count, pagenum int) int {
	var i, r int
	i, r = count / pagenum, count % pagenum
	if r == 0 {
		return i
	}
	return i + 1
}

// 计算分页起点和终点
// @params count 数据集长度
// @params pageNum 单页数据长度
// @params page    页码索引
// @defaultPageNum 单页数据最高长度
func Page(count, pageNum, page, defaultPageNum int) (int, int) {
	var start int
	if pageNum == 0 || pageNum > defaultPageNum {
		pageNum = defaultPageNum
	}
	// 页码大于1时需要计算是否超出范围
	if (page > 1) && (page > CalcPage(count, pageNum)) {
		page = 1
	}
	if page == 0 {
		page = 1
	}
	start = (page - 1) * pageNum
	if start > count {
		return 0, pageNum
	}
	// 避免超出
	if n := count - start; n < pageNum {
		pageNum = n
	}
	return start, pageNum
}

