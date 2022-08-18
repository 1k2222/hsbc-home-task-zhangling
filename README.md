# hsbc-home-task-zhangling

# Design

![](.//design.png)

From bottom to top, the system was designed in following three layers: 

1. **Entity layer**: The entity abstracted by the requirement. They are User, Role and AuthToken. Because the entity of this task is simple, all of them are designed by Anemic Domain Model.
2. **Repo layer**: Each element in this layer implements but only implements basic CURD logic of corresponding entity. In this layer, business logics such as hashing password and expiring auth tokens are excepted.
3. **Service layer**: All the APIs are implemented in this layer by invoking repo layer and implementing additional business logics. Constraints between repos are also implemented in this layer.

**Access layer**, not implemented in this task for concise, can be added above Service layer to support invocation by HTTP/GRPC.

# Test

Test file is service_test.go, which contains following group of test cases:

1. **TestUserRoleCreateAndDelete**: Test simple Create and Delete operation of User and Role. 
2. **TestAuthenticateAndInvalidate**: Test authenticate/invalidate operations as well as the validity of auth token after its corresponding user is deleted.
3. **TestTestTokenExpireTime**: Test whether token will expire after pre-configured time.
4. **TestUserRole**: Test complex operations of user/role relation.

To run all test cases, run the following command:

```console
$ go test service_test.go
```