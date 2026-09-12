# 测试用例与追踪

| Case ID | 需求 | 操作与预期 | 自动化位置 |
| --- | --- | --- | --- |
| DOC-01 | AC-01 | 注入新项目，入口指向 docs，未验证知识明确标记 | test_harness.FrameworkTest.test_install_into_non_go_project_is_idempotent_and_preserves_app |
| TASK-01 | AC-02 | 创建新任务，步骤材料齐全；再次创建同 ID 失败且文件不变 | test_tasks.TaskWorkflowTest.test_create_task_and_refuse_duplicate_without_losing_work |
| TASK-02 | AC-03 | 删除步骤文件或移除 status 登记，检查失败 | test_tasks.TaskWorkflowTest.test_missing_material_and_unregistered_task_are_detected |
| TASK-03 | AC-02/03 | UI 未判定不能 active，验证未通过不能 done | test_tasks.TaskWorkflowTest.test_start_requires_ui_applicability_and_completion_requires_verification |
| TASK-04 | AC-03 | 更新状态只修改任务表，保留表外历史记忆 | test_tasks.TaskWorkflowTest.test_status_update_keeps_memory_outside_the_table |
| TASK-05 | AC-02 | 页面用例缺少每步预期时不能通过格式检查 | test_tasks.TaskWorkflowTest.test_browser_case_requires_observable_step_assertions |
| TASK-06 | AC-02 | 非法 ID 或任务目录软链接不能写入目录外 | test_tasks.TaskWorkflowTest.test_invalid_task_id_or_symlink_cannot_write_outside_tasks |
| UPGRADE-01 | AC-06 | 升级保留项目文档、状态记忆和自写 API 测试 | test_tasks.TaskWorkflowTest.test_upgrade_preserves_documents_cases_status_and_api_tests |
| API-01 | AC-04 | 注册后登录，响应 code/status 为 200、Token 非空；带 Token 访问 hello 返回 Hello, API! | AuthAPITest.test_register_login_and_access_protected_hello |
| API-02 | AC-04 | 空/非法注册和登录参数为 400、错误密码和未知用户为 401、重复用户名为 409；响应为 JSON 且 message 非空 | AuthAPITest 的 invalid_registration、missing_login_fields、wrong_password_and_unknown_user、duplicate_username 用例 |
| API-03 | AC-04 | 缺少或无效 Token 访问 hello 为 401；不泄露 Token | AuthAPITest 的 protected_route_requires_token、protected_route_rejects_invalid_token 用例 |
| SKILL-01 | AC-05 | 代码变更后执行 skill，记录新增/复用用例和运行结果 | 本任务验证记录 |

页面操作集合见 [ui-cases.json](ui-cases.json)。本任务无前端页面，集合明确不适用。
API 与框架用例不会被计作浏览器执行结果。

自动化源：[框架回归](../../../.agents/framework/tests/test_harness.py)、
[任务回归](../../../.agents/framework/tests/test_tasks.py)、[API 用例](../../../tests/api/test_auth.py)。
前提/清理：框架测试使用临时项目并自动清理；API runner 每次创建独立内存服务，用例用户名互不影响，退出后关闭进程并删除临时产物。
