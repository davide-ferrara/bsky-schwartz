import argparse
import sys

from .api_client import APIClient
from .analyzer import Analyzer
from .charts import (
    plot_value_comparison,
    plot_value_differences,
    plot_cluster_differences,
)


def main():
    parser = argparse.ArgumentParser(
        description="Benchmark tool for comparing GPT vs Minimax Schwartz value analysis"
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
        default="comparison_chart.png",
        help="Output PNG file path",
    )
    parser.add_argument(
        "--diff",
        action="store_true",
        help="Show difference chart instead of side-by-side comparison",
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

    print("Analyzing with Minimax...")
    minimax_results = [analyzer.analyze_uri(uri, "minimax") for uri in uris]
    minimax_results = [r for r in minimax_results if r is not None]

    if not gpt_results or not minimax_results:
        print("Error: No results from one or both models")
        sys.exit(1)

    results = {"gpt": gpt_results, "minimax": minimax_results}

    if args.diff:
        diffs = analyzer.compute_differences(results)
        print("\nValue Differences (GPT - Minimax):")
        for val, diff in sorted(diffs.items(), key=lambda x: abs(x[1]), reverse=True):
            sign = "+" if diff >= 0 else ""
            print(f"  {val}: {sign}{diff:.2f}")

        if args.clusters:
            cluster_diffs = analyzer.compute_cluster_diffs(diffs)
            plot_cluster_differences(
                cluster_diffs,
                output_path=args.output,
                title="GPT vs Minimax: Cluster Differences",
            )
        else:
            plot_value_differences(
                diffs,
                output_path=args.output,
                title="GPT vs Minimax: Value Differences",
            )
    else:
        print("\nAverage Scores:")
        values_list = [
            "sd_thought",
            "sd_action",
            "stimulation",
            "hedonism",
            "achievement",
            "dominance",
            "resources",
            "face",
            "personal_sec",
            "societal_sec",
            "tradition",
            "rule_conf",
            "inter_conf",
            "humility",
            "caring",
            "dependability",
            "universalism",
            "nature",
            "tolerance",
        ]
        for val in values_list:
            gpt_scores = [r.values.get(val, 0) for r in gpt_results]
            mm_scores = [r.values.get(val, 0) for r in minimax_results]
            gpt_avg = sum(gpt_scores) / len(gpt_scores) if gpt_scores else 0
            mm_avg = sum(mm_scores) / len(mm_scores) if mm_scores else 0
            print(f"  {val}: GPT={gpt_avg:.2f}, Minimax={mm_avg:.2f}")

        plot_value_comparison(
            results,
            output_path=args.output,
            title="GPT vs Minimax: Value Comparison",
        )

    print(f"\nResults saved to {args.output}")


if __name__ == "__main__":
    main()
