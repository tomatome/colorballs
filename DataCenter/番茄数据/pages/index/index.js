//index.js
//获取应用实例
var app = getApp()
Page({
  data: {
    userInfo: {},
    ballslist: []
  },
  //事件处理函数
  bindViewTap: function() {
    wx.navigateTo({
      url: '../logs/logs'
    })
  },
  onLoad: function () {
    console.log("this is onLoad")
    var that = this
  	//调用应用实例的方法获取全局数据
    app.getUserInfo(function(userInfo){
      //更新数据
      that.setData({
        userInfo:userInfo
      })
      that.update()
    })
    wx.request({  
      url: 'https://golang-mework.rhcloud.com', 
      data: {},  
      method: 'POST', 
      //header: {},
      success: function(res){
        console.log("this is request")
        console.log(res.statusCode)
        that.setData({
          ballslist:res.data  
        })  
      },  
      fail: function() {  
        // fail  
      },  
      complete: function() {  
        // complete  
      }  
    })
  }
})
