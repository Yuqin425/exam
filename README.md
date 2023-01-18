

# 蓝山考核

### 代码架构

```
├── README.md           					// 说明文档
	├── conf
		└── config.ymal						// 配置文件
    ├── api									// 接口层
    │   └── middlewares						// 中间件
    		└── jwt
    ├── service								// 业务逻辑层
    ├── dao									// 数据库层
 		└──	mysql
    ├── models								// 模型层
    ├── utils
    ├── settings
    ├── logger
    ├── log									// 项目日志
    ├── go.mod
    └── main.go
```


## 功能

### 用户注册登录

**首先当然是建立模型**

`model/user.go`

```go
type User struct {
	UserID       uint64 `json:"user_id,string" db:"user_id"`
	UserName     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	AccessToken  string
	RefreshToken string
}

type RegisterForm struct {
	UserName        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginForm struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
```



**注册业务逻辑及数据库操作**

1. 判断用户是否存在
2. 应用雪花算法生成唯一id
3. 新增用户数据

`service/user.go`

```go
func SignUp(p *models.RegisterForm) (error error) {
	// 1、判断用户存不存在
	err := mysql.CheckUserExist(p.UserName)
	if err != nil {
		// 数据库查询出错
		return err
	}

	// 2、生成UID
	userId, err := snowflake.GetID()
	if err != nil {
		return mysql.ErrorGenIDFailed
	}
	// 构造一个User实例
	u := models.User{
		UserID:   userId,
		UserName: p.UserName,
		Password: p.Password,
	}
	// 3、保存进数据库
	return mysql.InsertUser(&u)
}
```

`dao/mysql/user.go`

```go
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

// 检验用户名是否存在
func CheckUserExist(username string) (error error) {
	sqlstr := `select count(*) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExit
	}
	return
}

// 插入用户数据
func InsertUser(user *models.User) (error error) {
	// 对密码进行加密
	user.Password = encryptPassword([]byte(user.Password))
	// 执行SQL语句入库
	sqlstr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err := db.Exec(sqlstr, user.UserID, user.UserName, user.Password)
	return err
}
```

**登录业务逻辑及数据库操作**

1. 数据库能否查找到用户名，不能知道则返回。
2. 验证密码是否正确。
3. 正确则返回我们的 token.

`service/user.go`

```go
func Login(p *models.LoginForm) (user *models.User, error error) {
	user = &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT
	//return jwt.GenToken(user.UserID,user.UserName)
	atoken, rtoken, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.AccessToken = atoken
	user.RefreshToken = rtoken
	return
}
```

`dao/mysql/user.go`

```go
func Login(user *models.User) (err error) {
	originPassword := user.Password // 记录一下原始密码(用户登录的密码)
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		return ErrorUserNotExit
	}
	// 生成加密密码与查询到的密码比较
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return
}
```

**开始编写api层**

*从前端获取数据并进行业务操作*

`api/user.go`

```go
// 注册
func SignUpHandler(c *gin.Context) {
	// 1.获取请求参数 2.校验数据有效性
	var fo *models.RegisterForm
	if err := c.ShouldBindJSON(&fo); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			utils.ResponseError(c, utils.CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 3.业务处理——注册用户
	if err := service.SignUp(fo); err != nil {
		zap.L().Error("service.signup failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExit) {
			utils.ResponseError(c, utils.CodeUserExist)
			return
		}
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	//返回响应
	utils.ResponseSuccess(c, nil)
}

// 登录
func LoginHandler(c *gin.Context) {
	// 获取请求参数及参数校验
	var u *models.LoginForm
	if err := c.ShouldBindJSON(&u); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			utils.ResponseError(c, utils.CodeInvalidParams) // 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2、业务逻辑处理——登录
	user, err := service.Login(u)
	if err != nil {
		zap.L().Error("service.Login failed", zap.String("username", u.UserName), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExit) {
			utils.ResponseError(c, utils.CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorPasswordWrong) {
			utils.ResponseError(c, utils.CodeInvalidPassword)
			return
		}
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 3、返回响应
	utils.ResponseSuccess(c, gin.H{
		"user_id":       fmt.Sprintf("%d", user.UserID),
		"user_name":     user.UserName,
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}
```

到目前为止，注册登录的功能就大致完成了



## 关于题目的相关功能

**还是先建立模型**

```go
type Problem struct {
	ProblemID   uint64    `json:"problem_id,string" db:"problem_id"`
	AuthorId    uint64    `json:"author_id" db:"author_id"`
	CommunityID uint64    `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	Input       string    `json:"input" db:"input" binding:"required"`
	Output      string    `json:"output" db:"output" binding:"required"`
	CreateTime  time.Time `json:"-" db:"create_time"`
}
```

发布题目与上面注册用户的实现方法大差不差，我这里也不重复说了，这里主要是获取问题和修改删除问题的实现

详情看源码 

**获取问题业务逻辑实现**

`service/problem.go`

获取问题大致逻辑非常简单，就是通过问题id直接查找问题详情，列表直接查询

```go
// 根据问题id查询题目详情
func GetProblemById(problemID int64) (data *models.ApiProblemDetail, err error) {
	// 查询信息
	problem, err := mysql.GetProblemByID(problemID)
	if err != nil {
		zap.L().Error("mysql.GetProblemByID(problemID) failed",
			zap.Int64("problemID", problemID),
			zap.Error(err))
		return nil, err
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(problem.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed",
			zap.Uint64("AuthorID", problem.AuthorId),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityByID(problem.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Uint64("community_id", problem.CommunityID),
			zap.Error(err))
		return
	}
	// 接口数据拼接
	data = &models.ApiProblemDetail{
		Problem:         problem,
		CommunityDetail: community,
		AuthorName:      user.UserName,
	}
	return
}

// 获取所有题目列表
func GetProblemList(page, size int64) (data []*models.ApiProblemDetail, err error) {
	problemList, err := mysql.GetProblemList(page, size)
	if err != nil {
		fmt.Println(err)
		return
	}
	data = make([]*models.ApiProblemDetail, 0, len(problemList)) // data 初始化
	for _, problem := range problemList {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(problem.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("problemID", problem.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityByID(problem.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", problem.CommunityID),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		problemdetail := &models.ApiProblemDetail{
			Problem:         problem,
			CommunityDetail: community,
			AuthorName:      user.UserName,
		}
		data = append(data, problemdetail)
	}
	return
}
```

**修改问题逻辑实现**

1. 查询是否登录
2. 当前用户是否为作者
3. 是作者直接根据id修改
4. 比较内容是否不同对数据库更新

`api/problem.go`

```go
func ProblemUpdateHandler(c *gin.Context) {
	// 获取参数及校验参数
	var newProblem models.Problem
	if err := c.ShouldBindJSON(&newProblem); err != nil {
		zap.L().Debug("c.ShouldBindJSON(problem) err", zap.Any("err", err))
		zap.L().Error("create problem with invalid param")
		utils.ResponseErrorWithMsg(c, utils.CodeInvalidParams, err.Error())
		return
	}

	// 获取参数(从URL中获取id)
	problemIdStr := c.Param("id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	pastProblem, err := service.GetProblemById(problemId)
	if err != nil {
		zap.L().Error("get problem detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
	}

	// 获取作者ID
	UserID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	ok := UserID == pastProblem.AuthorId
	if !ok {
		zap.L().Error("update problem with invalid param")
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}

	problem, err := service.UpdateProblem(&newProblem, problemId)
	if err != nil {
		zap.L().Error("service.UpdateProblem() failed")
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	// 3、返回响应
	utils.ResponseSuccess(c, problem)
}
```

**删除问题逻辑实现**

1. 查询是否登录
2. 当前用户是否为作者
3. 是作者直接根据id删除

`api/problem.go`

```go
func ProblemDeleteHandler(c *gin.Context) {
	// 获取参数(从URL中获取id)
	problemIdStr := c.Param("id")
	problemId, err := strconv.ParseInt(problemIdStr, 10, 64)
	problem, err := service.GetProblemById(problemId)
	if err != nil {
		zap.L().Error("get problem detail with invalid param", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}

	// 获取作者ID
	UserID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	ok := UserID == problem.AuthorId
	if !ok {
		zap.L().Error("delete problem with invalid param")
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}
	service.DeleteProblem(problemId)

	utils.ResponseSuccess(c, nil)
}
```



## 关于代码评测的功能

QAQ由于本人太菜了，实在无法完成



接下来就完善一下其他功能吧

## 发布题解功能

有人发问题也应该有人发答案

**先建立模型**

`models/answer.go`

```go
type Answer struct {
	ProblemID  uint64    `json:"problem_id,string" db:"problem_id" binding:"required"`
	ParentID   uint64    `db:"parent_id" json:"parent_id"`
	AnswerID   uint64    `db:"answer_id" json:"answer_id"`
	AuthorID   uint64    `json:"author_id" db:"author_id"`
	Content    string    `json:"content" db:"content" binding:"required"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
}
```

发布题解实现逻辑和发布问题差不多

1. 得到题目id
2. 雪花算法生成唯一id
3. 获取当前用户id并检查是否登录
4. 插入数据库

`api/answer.go`

```go
func AnswerHandler(c *gin.Context) {
	var answer models.Answer
	if err := c.BindJSON(&answer); err != nil {
		fmt.Println(err)
		utils.ResponseError(c, utils.CodeInvalidParams)
		return
	}
	// 生成ID
	answerID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	// 获取作者ID，当前请求的UserID
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeNotLogin)
		return
	}
	answer.AnswerID = answerID
	answer.AuthorID = userID

	// 创建题解
	if err := mysql.CreateAnswer(&answer); err != nil {
		zap.L().Error("mysql.CreateAnswer(&answer) failed", zap.Error(err))
		utils.ResponseError(c, utils.CodeServerBusy)
		return
	}
	utils.ResponseSuccess(c, nil)
}
```

剩下的列表功能和获取详情几乎就是和题目功能差不多了，基本是cv的，我这里就不啰嗦了，具体看原码
