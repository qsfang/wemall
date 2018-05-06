package cart

import (
	"encoding/json"
	"fmt"
	"gopkg.in/kataras/iris.v6"
	"strconv"
	"wemall/config"
	"wemall/controller/common"
	"wemall/model"
)

// Create 购物车中添加商品
func Create(ctx *iris.Context) {
	SendErrJSON := common.SendErrJSON
	var cart model.Cart

	if ctx.ReadJSON(&cart) != nil {
		SendErrJSON("参数错误", ctx)
		return
	}

	if cart.Count <= 0 {
		SendErrJSON("count不能小于0", ctx)
		return
	}

	var product model.Product
	if model.DB.First(&product, cart.ProductID).Error != nil {
		SendErrJSON("错误的商品id", ctx)
		return
	}

	session := ctx.Session()
	openID := session.GetString("weAppOpenID")

	if openID == "" {
		SendErrJSON("登录超时", ctx)
		return
	}

	cart.OpenID = openID
	if model.DB.Create(&cart).Error != nil {
		SendErrJSON("error", ctx)
		return
	}
	ctx.JSON(iris.StatusOK, iris.Map{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"id": cart.ID,
		},
	})
	return
}

//add struct Cart
type CartProduct struct {
	Cart    model.Cart    `json:"cart"`
	Product model.Product `json:"product"`
}

// List 商品列表
func List(ctx *iris.Context) {
	SendErrJSON := common.SendErrJSON
	var carts []model.Cart
	var product model.Product
	var cartProducts []CartProduct

	session := ctx.Session()
	openID := session.GetString("weAppOpenID")

	if openID == "" {
		SendErrJSON("登录超时", ctx)
		return
	}

	pageNo, err := strconv.Atoi(ctx.FormValue("pageNo"))

	if err != nil || pageNo < 1 {
		pageNo = 1
	}

	offset := (pageNo - 1) * config.ServerConfig.PageSize

	//默认按创建时间，降序来排序
	var orderStr = "created_at"
	if ctx.FormValue("asc") == "1" {
		orderStr += " asc"
	} else {
		orderStr += " desc"
	}

	pageSize := config.ServerConfig.PageSize
	queryErr := model.DB.Offset(offset).Limit(pageSize).Order(orderStr).Where("open_id = ?", openID).Find(&carts).Error

	if queryErr != nil {
		SendErrJSON("error", ctx)
		return
	}

	for _, cart := range carts {
		err := model.DB.First(&product, cart.ProductID).Error
		if err != nil {
			fmt.Println(err.Error())
			SendErrJSON("error", ctx)
			return
		}

		if model.DB.First(&product.Image, product.ImageID).Error != nil {
			product.Image = model.Image{}
		}

		var imagesSQL []uint
		if err := json.Unmarshal([]byte(product.ImageIDs), &imagesSQL); err == nil {
			var images []model.Image
			if model.DB.Where("id in (?)", imagesSQL).Find(&images).Error != nil {
				product.Images = nil
			} else {
				product.Images = images
			}
		} else {
			product.Images = nil
		}

		if err := model.DB.Model(&product).Related(&product.Categories, "categories").Error; err != nil {
			fmt.Println(err.Error())
			SendErrJSON("error", ctx)
			return
		}

		if product.HasProperty {
			if err := model.DB.Model(&product).Related(&product.Properties).Error; err != nil {
				fmt.Println(err.Error())
				SendErrJSON("error", ctx)
				return
			}

			for i := 0; i < len(product.Properties); i++ {
				property := product.Properties[i]
				if err := model.DB.Model(&property).Related(&property.PropertyValues).Error; err != nil {
					fmt.Println(err.Error())
					SendErrJSON("error", ctx)
					return
				}
				product.Properties[i] = property
			}

			if err := model.DB.Model(&product).Related(&product.Inventories).Error; err != nil {
				fmt.Println(err.Error())
				SendErrJSON("error", ctx)
				return
			}

			for i := 0; i < len(product.Inventories); i++ {
				inventory := product.Inventories[i]
				if err := model.DB.Model(&inventory).Related(&inventory.PropertyValues, "property_values").Error; err != nil {
					fmt.Println(err.Error())
					SendErrJSON("error", ctx)
					return
				}
				product.Inventories[i] = inventory
			}
		}
		cp := CartProduct{
			Cart:    cart,
			Product: product,
		}

		cartProducts = append(cartProducts, cp)
	}

	ctx.JSON(iris.StatusOK, iris.Map{
		"errNo": model.ErrorCode.SUCCESS,
		"msg":   "success",
		"data": iris.Map{
			"carts": cartProducts,
		},
	})
}
