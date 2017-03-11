resource "aws_elasticache_cluster" "ptt-redis" {
    cluster_id = "ptt-redis"
    engine = "redis"
    node_type = "cache.t2.micro"
    port = 6379
    num_cache_nodes = 1
    parameter_group_name = "default.redis3.2"
    security_group_ids = ["${aws_security_group.ecs.id}"]
}
