###simulate（单元测试模块 + 自动化测试） 
 - 单元测试 U_xxx 包 内置可自行实现golang 单元测试模块；
 - 自动化测试如下：
 		
		文件名以“T”+"_" + 模块名
        json配置格式（支持复杂protobuf消息结构）：
		[{
	        "connAddr":"127.0.0.1:51001",
	        "module":"game",
	        "route": 8,
	        "mainId": 4, 
	        "msg":[
	            {
	                "request":{
	                    "subId": 1, 
	                    "name": "CS_EnterServer_Req", 
	                    "params":{
	                       "A": 1,
	                       "B":{
	                           "B1":1,
	                           "B2":"2"
	                       },
	                       "C":{
	                            "C1":{
	                                "C11":1,
	                                "C12":"12"
	                            }
	                       }
	                    }
	                }, 
	
	                "response":{
	                    "subId": 2, 
	                    "name": "SC_EnterServer_Rsp", 
	                    "params":{
	
	                    }
	                }
	            }
	       ]
	    }]

		连接gateway网关json配置：
		[{"ConnAddr": "127.0.0.1:51001"}]

	自加载完所有配置后，依次自动向服务器发送消息测试进行验证.