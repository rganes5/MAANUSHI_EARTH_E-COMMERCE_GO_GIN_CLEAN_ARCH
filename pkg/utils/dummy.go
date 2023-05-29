package utils

// UserSignup
// @Summary api for Signup a new user
// @ID Signup-user
// @Description Create a new user with the specified details.
// @Tags Users Signup
// @Accept json
// @Produce json
// @Param user_details body model.NewUserInfo true "New user Details"
// @Success 200 {object} res.Response
// @Failure 500 {object} res.Response
// @Failure 422 {object} res.Response
// @Router /user/signup [post]

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

/////////////////////////////////////

// func (cr *AdminHandler) SignUp(c *gin.Context) {
// 	var admin domain.Admin
// 	if err := c.BindJSON(&admin); err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}

// 	if ok := support.Email_validater(admin.Email); !ok {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": "Email format incorrect",
// 		})
// 		return
// 	}

// 	if ok := support.MobileNum_validater(admin.MobileNum); !ok {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
// 			"error": "Not a valid mobile number",
// 		})
// 		return
// 	}
// 	if _, err := cr.AdminUseCase.FindbyEmail(c.Request.Context(), admin.Email); err == nil {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"error": "User already Exsists",
// 		})
// 		return
// 	}

// 	admin.Password, _ = support.HashPassword(admin.Password)
// 	err := cr.AdminUseCase.SignUpAdmin(c.Request.Context(), admin)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"User registration": "Success",
// 	})
// }
