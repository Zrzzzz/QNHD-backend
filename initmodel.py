modelsmallname = input('模板全小写英文: ')
modelbigname = input('模板大驼峰英文: ')
modelcnname = input('模板中文名: ')

template = """package v1

import "github.com/gin-gonic/gin"

func Get{0}(c *gin.Context) {{

}}

func Add{0}(c *gin.Context) {{

}}

func Edit{0}(c *gin.Context) {{
	
}}

func Delete{0}(c *gin.Context) {{

}}
""".format(modelbigname)

inittmp = """func init{1}(g *gin.RouterGroup) {{
	//获取{2}列表
	g.GET("/{0}", v1.Get{1})
	//新建{2}
	g.POST("/{0}", v1.Add{1})
	//修改{2}
	g.PUT("/{0}/:id", v1.Edit{1})
	//删除指定{2}
	g.DELETE("/{0}/:id", v1.Delete{1})
}}
""".format(modelsmallname, modelbigname, modelcnname)


with open('{}.go'.format(modelsmallname), 'w') as f:
	f.write(template)
	
print(inittmp)
	

