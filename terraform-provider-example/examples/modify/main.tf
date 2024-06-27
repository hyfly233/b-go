resource "example_modifier" "router" {
  name                 = "test_router_12"
  tenant_id            = "5a6275f8e2794214b464598f9086b5d7"
  description          = "test_router"
  external_network_id  = "ba70d9f0-655b-4520-be21-a9bac82e305d"
  external_enable_snat = false
  external_fixed_ips = [
    {
      subnet_id = "a15a53f3-85b2-4e07-83b9-8d7a1a672451"
    }
  ]
}
