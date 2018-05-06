var config   = require('../../config/config.js');

Page({
    data: {
        appName: config.name,
        totalMoney: 0,
        carts: []
    },
    onLoad: function() {
      var self = this;
      var app = getApp();
      var totalMoney = 0;
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
          var totalMoney = 0;
          console.log("Get user carts", carts)

          // request cart product info
          for (var i = 0; i < carts.length; i++) {
            var cp = carts[i]
            cp.product.image.url = config.static.imageDomain + cp.product.image.url;
            totalMoney += cp.cart.count * cp.product.price;
          }
          self.setData({
            carts: carts,
            totalMoney: totalMoney,
          }); 
      }
    });
  }
})