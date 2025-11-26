#!/usr/bin/env python3
"""
Generate AWS architecture diagram for Gestione Caselo infrastructure.

Requirements:
    pip install diagrams

Usage:
    python generate_diagram.py

Output:
    gestione_caselo_architecture.png
"""

from diagrams import Diagram, Cluster, Edge
from diagrams.aws.network import Route53, CloudFront, APIGateway
from diagrams.aws.storage import S3
from diagrams.aws.compute import Lambda
from diagrams.aws.database import Dynamodb
from diagrams.aws.integration import SQS
from diagrams.aws.engagement import SES
from diagrams.aws.management import Cloudwatch
from diagrams.aws.security import IAM, IdentityAndAccessManagementIam
from diagrams.onprem.vcs import Github
from diagrams.onprem.client import Users

# Diagram configuration
graph_attr = {
    "fontsize": "14",
    "bgcolor": "white",
    "pad": "0.5",
    "splines": "ortho",
}

node_attr = {
    "fontsize": "12",
}

edge_attr = {
    "fontsize": "10",
}

with Diagram(
    "Gestione Caselo - AWS Architecture",
    filename="gestione_caselo_architecture",
    show=False,
    direction="TB",
    graph_attr=graph_attr,
    node_attr=node_attr,
    edge_attr=edge_attr,
    outformat="png",
):
    users = Users("Users")

    with Cluster("DNS & CDN"):
        dns = Route53("Route53\n*.gestionecaselo.it")
        cdn = CloudFront("CloudFront\nCDN + SSL")

    with Cluster("Frontend"):
        frontend_bucket = S3("S3 Bucket\nVue.js SPA")

    with Cluster("API Layer"):
        api = APIGateway("API Gateway\nHTTP API")

    with Cluster("Compute"):
        eventform_lambda = Lambda("EventForm\nGraphQL API")
        emails_lambda = Lambda("Emails\nEmail Sender")

    with Cluster("Data Storage"):
        dynamodb = Dynamodb("DynamoDB\nSingle Table")

    with Cluster("Message Queues"):
        emails_queue = SQS("Emails Queue")
        emails_dlq = SQS("Emails DLQ")

    with Cluster("Email Delivery"):
        ses = SES("SES")

    with Cluster("Observability"):
        logs = Cloudwatch("CloudWatch\nLogs & Alarms")

    with Cluster("CI/CD & Global"):
        github = Github("GitHub Actions")
        oidc = IdentityAndAccessManagementIam("OIDC Provider")
        iam = IAM("IAM Roles")
        artifacts = S3("Lambda Artifacts")

    # User flow - Frontend
    users >> Edge(label="HTTPS") >> dns
    dns >> Edge(label="routes") >> cdn
    cdn >> Edge(label="serves") >> frontend_bucket

    # User flow - API
    users >> Edge(label="GraphQL") >> api
    api >> Edge(label="invoke") >> eventform_lambda

    # EventForm Lambda interactions
    eventform_lambda >> Edge(label="read/write") >> dynamodb
    eventform_lambda >> Edge(label="enqueue") >> emails_queue

    # Email processing flow
    emails_queue >> Edge(label="trigger") >> emails_lambda
    emails_queue >> Edge(label="failed", style="dashed") >> emails_dlq
    emails_lambda >> Edge(label="send") >> ses

    # Logging
    eventform_lambda >> Edge(style="dotted") >> logs
    emails_lambda >> Edge(style="dotted") >> logs

    # CI/CD flow
    github >> Edge(label="authenticate") >> oidc
    oidc >> Edge(label="assume") >> iam
    iam >> Edge(label="deploy", style="dashed") >> eventform_lambda
    iam >> Edge(label="deploy", style="dashed") >> emails_lambda
    iam >> Edge(label="upload", style="dashed") >> frontend_bucket
    iam >> Edge(label="store", style="dashed") >> artifacts

print("✓ Diagram generated: gestione_caselo_architecture.png")
