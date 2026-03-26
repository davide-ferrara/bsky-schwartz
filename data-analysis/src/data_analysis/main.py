import argparse
import sys

from .api_client import APIClient
from .analyzer import Analyzer
from .charts import plot_value_differences, plot_cluster_differences


def main():
    parser = argparse.ArgumentParser(
        description="Benchmark tool for comparing GPT vs Qwen Schwartz value analysis"
    )
    parser.add_argument(
        "--base-url",
        default="http://localhost:8080",
        help="Base URL of the bsky-schwartz API",
    )
    parser.add_argument(
        "--query",
        help="Search query to get URIs (alternative to --uris)",
    )
    parser.add_argument(
        "--limit",
        type=int,
        default=10,
        help="Number of posts to fetch from query",
    )
    parser.add_argument(
        "--uris",
        help="Comma-separated list of URIs (alternative to --query)",
    )
    parser.add_argument(
        "--output",
        default="diff_chart.png",
        help="Output PNG file path",
    )
    parser.add_argument(
        "--clusters",
        action="store_true",
        help="Show cluster-level differences instead of individual values",
    )

    args = parser.parse_args()

    api = APIClient(args.base_url)

    if not api.health_check():
        print("Error: API health check failed. Is the server running?")
        sys.exit(1)

    if args.uris:
        uris = [u.strip() for u in args.uris.split(",")]
    elif args.query:
        print(f"Fetching URIs for query: {args.query}")
        uris = api.search_uris(args.query, args.limit)
        print(f"Found {len(uris)} URIs")
    else:
        print("Error: Provide either --query or --uris")
        sys.exit(1)

    analyzer = Analyzer(api)

    print("Analyzing with GPT...")
    gpt_results = [analyzer.analyze_uri(uri, "gpt") for uri in uris]
    gpt_results = [r for r in gpt_results if r is not None]

    print("Analyzing with Qwen...")
    qwen_results = [analyzer.analyze_uri(uri, "qwen") for uri in uris]
    qwen_results = [r for r in qwen_results if r is not None]

    if not gpt_results or not qwen_results:
        print("Error: No results from one or both models")
        sys.exit(1)

    results = {"gpt": gpt_results, "qwen": qwen_results}

    diffs = analyzer.compute_differences(results)

    print("\nValue Differences (GPT - Qwen):")
    for val, diff in sorted(diffs.items(), key=lambda x: abs(x[1]), reverse=True):
        sign = "+" if diff >= 0 else ""
        print(f"  {val}: {sign}{diff:.2f}")

    if args.clusters:
        cluster_diffs = analyzer.compute_cluster_diffs(diffs)
        plot_cluster_differences(
            cluster_diffs,
            output_path=args.output,
            title="GPT vs Qwen: Cluster Differences",
        )
    else:
        plot_value_differences(
            diffs,
            output_path=args.output,
            title="GPT vs Qwen: Value Differences",
        )

    print(f"\nResults saved to {args.output}")


if __name__ == "__main__":
    main()
