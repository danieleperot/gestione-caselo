output "hosted_zone_id" {
  value       = aws_route53_zone.main.zone_id
  description = "ID of the Route53 hosted zone created for the domain"
}

output "ns_records" {
  value       = aws_route53_zone.main.name_servers
  description = "List of nameserver records that must be configured on the parent domain to point the subdomain of the app to the Route53 hosted zone"

}
