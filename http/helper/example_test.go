package helper

import (
	"fmt"
	"github.com/ArtisanCloud/PowerLibs/v2/http/contract"
	"log"
	"net/http"
)

func ExampleRequestHelper_WithMiddleware() {
	// 初始化 helper
	helper, err := NewRequestHelper(&Config{
		BaseUrl: "https://www.baidu.com",
	})
	if err != nil {
		log.Fatalln(err)
	}

	helper.WithMiddleware(func(handle contract.RequestHandle) contract.RequestHandle {
		return func(request *http.Request) (response *http.Response, err error) {
			// 前置中间件
			fmt.Println("这里是前置中间件1, 在请求前执行:")
			fmt.Println(request.URL.String())

			response, err = handle(request)
			// handle 执行之后就可以操作 response 和 err

			// 后置中间件
			fmt.Println("这里是后置置中间件1, 在请求后执行:")
			if err == nil {
				fmt.Println(request.URL.String())
				fmt.Println(response.Status)
			}
			return
		}
	})

	logMiddleware := func(logger *log.Logger) contract.RequestMiddleware {
		return contract.RequestMiddleware(func(handle contract.RequestHandle) contract.RequestHandle {
			return func(request *http.Request) (response *http.Response, err error) {
				// 前置中间件
				logger.Println("这里是前置中间件log, 在请求前执行")

				response, err = handle(request)
				// handle 执行之后就可以操作 response 和 err

				// 后置中间件
				logger.Println("这里是后置置中间件log, 在请求后执行")
				return
			}
		})
	}

	helper.WithMiddleware(func(handle contract.RequestHandle) contract.RequestHandle {
		return func(request *http.Request) (response *http.Response, err error) {
			// 前置中间件
			fmt.Println("这里是前置中间件2, 在请求前执行")

			response, err = handle(request)
			// handle 执行之后就可以操作 response 和 err

			// 后置中间件
			fmt.Println("这里是后置置中间件2, 在请求后执行")
			return
		}
	}, logMiddleware(log.Default()))

	resp, err := helper.Df().Method(http.MethodGet).Request()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(resp.Request.URL.String())
	log.Println(resp.Status)
	// Output:
	// 这里是前置中间件1, 在请求前执行:
	// https://www.baidu.com
	// 这里是前置中间件2, 在请求前执行
	// 这里是后置置中间件2, 在请求后执行
	// 这里是后置置中间件1, 在请求后执行:
	// https://www.baidu.com
	// 200 OK
}

func ExampleRequestHelper_Df() {
	// 初始化 helper
	helper, err := NewRequestHelper(&Config{})
	if err != nil {
		log.Fatalln(err)
	}

	var result struct {
		Status string
	}

	// mock server response: {"status":"success"}
	err = helper.Df().Method(http.MethodGet).
		Url("http://localhost:3000/do-testing").
		Header("a", "b").
		Query("a[]", "b", "c").
		Json(map[string]interface{}{
			"a": "b",
		}).
		Result(&result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
	// Output: {success}
}
