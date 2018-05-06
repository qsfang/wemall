var config   = require('../../config/config.js');

Page({
    data: {
        appName: config.name,
        carts: []
    },
    onLoad: function() {
      var self = this
      var app = getApp();
      // todo 验证是否登陆
      console.log(app)
      wx.request({
        url: config.api.reqCartList,
        method: "GET",
        data: {
        },
        header: {
          'content-type': 'application/json',
          'Cookie': 'sid=' + app.globalData.sid
        },
        // request user carts
        success: function (res) {
          var carts = res.data.data.carts;
          var cartsInfo = []
          console.log("Get user carts", carts)

          // request cart product info
          for (var i = 0; i < carts.length; i++) {
            var productId = carts[i].productId
            var productCount = carts[i].count;
            wx.request({
              url: config.api.reqProductDetail.replace(':id', productId),
              success: function (res) {
                console.log("Cart product info response", res)
                var product = res.data.data.product || null;
                product.image.url = config.static.imageDomain + product.image.url;
                product.price = product.price.toFixed(2);
                var cartDetail = {}
                cartDetail.id = productId;
                cartDetail.count = productCount
                cartDetail.imageUrl = product.image.url;
                cartDetail.name = product.name;
                cartDetail.price = product.price;
                cartDetail.checked = true;
                cartsInfo.push(cartDetail);

                console.log("Carts Info: ", cartsInfo)
                self.setData({
                  carts: cartsInfo,
                }); 
            }
          }); 
        }
      }
    });
  }
})