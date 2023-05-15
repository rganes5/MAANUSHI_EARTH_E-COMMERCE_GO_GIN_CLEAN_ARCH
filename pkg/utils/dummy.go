package utils

//BLOCK AND UNBLOCK
//
// if err := Db.Where("id = ?", input.ID).First(&user).Error; err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
// 	return
// }

// block := !user.Block

// if err := Db.Model(&user).Where("id=?", input.ID).Update("block", block).Error; err != nil {
// 	c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to block"})
// 	return
// }
//
//lIST WITH PAGINATION
//
// func (pc *ProductController) ListProductsPagination(c *gin.Context) {

// // page := c.DefaultQuery("page", "1")
// // limit := c.DefaultQuery("limit", "10")
// page := c.Param("page")
// limit := c.Param("limit")

// pageNum, err := strconv.Atoi(page)
// if err != nil || pageNum < 1 {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
// 	return
// }

// limitNum, err := strconv.Atoi(limit)
// if err != nil || limitNum < 1 || limitNum > 100 {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
// 	return
// }

// offset := (pageNum - 1) * limitNum

// products, totalCount, err := pc.productService.ListPaginated(offset, limitNum)
// if err != nil {
// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 	return
// }

// response := gin.H{
// 	"page":          pageNum,
// 	"limit":         limitNum,
// 	"total_records": totalCount,
// 	"products":      products,
// }

// c.JSON(http.StatusOK, response)
// }
