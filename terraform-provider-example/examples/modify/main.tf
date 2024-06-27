resource "example_modifier" "mod" {
  replace             = "test_router_12"
  replace_if_configured = "5a6275f8e2794214b464598f9086b5d7"
  use_state_for_unknown  = "test_router"
  list_optional = [
    "test1", "test2"
  ]
}
