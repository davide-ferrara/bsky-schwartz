import json
import os
import glob
from charts import plot_values_comparison
from spider import plot_spider_chart


def find_latest_json(pattern: str) -> str:
    """Find the most recent JSON file matching the pattern."""
    files = glob.glob(pattern)
    if not files:
        return None
    # Sort by modification time, most recent first
    files.sort(key=os.path.getmtime, reverse=True)
    return files[0]


def load_json(file_path: str):
    if not file_path or not file_path.endswith("json"):
        return None

    with open(file_path) as f:
        data = json.load(f)

    return data


def calculate_values_avg(posts):
    if not posts:
        return {}

    values = posts[0]["ValueAnalysis"]["Rating"].copy()
    for post in posts:
        curr_values = post["ValueAnalysis"]["Rating"]
        for v in curr_values:
            values[v] += curr_values[v]
            values[v] /= 2

    return values


def main():
    # Trova automaticamente i file JSON più recenti per ogni modello
    gpt_file = find_latest_json("../post_gpt-4.1-mini*.json")
    mistral_file = find_latest_json("../post_ministral-14b*.json")
    deepseek_file = find_latest_json("../post_DeepSeek-V3*.json")
    qwen_file = find_latest_json("../post_Qwen3-VL*.json")

    print(f"GPT file: {gpt_file}")
    print(f"Mistral file: {mistral_file}")
    print(f"DeepSeek file: {deepseek_file}")
    print(f"Qwen file: {qwen_file}")

    # Carica i dati
    gpt_posts = load_json(gpt_file)
    mistral_posts = load_json(mistral_file)
    deepseek_posts = load_json(deepseek_file)
    qwen_posts = load_json(qwen_file)

    if not gpt_posts or not mistral_posts or not deepseek_posts or not qwen_posts:
        print("Error loading data")
        if not gpt_posts:
            print("  - GPT posts not found")
        if not mistral_posts:
            print("  - Mistral posts not found")
        if not deepseek_posts:
            print("  - DeepSeek posts not found")
        if not qwen_posts:
            print("  - Qwen posts not found")
        return

    print("\nData loaded successfully")
    print(f"GPT posts: {len(gpt_posts)}")
    print(f"Mistral posts: {len(mistral_posts)}")
    print(f"DeepSeek posts: {len(deepseek_posts)}")
    print(f"Qwen posts: {len(qwen_posts)}")

    # Calcola le medie per ogni modello
    gpt_avg_values = calculate_values_avg(gpt_posts)
    mistral_avg_values = calculate_values_avg(mistral_posts)
    deepseek_avg_values = calculate_values_avg(deepseek_posts)
    qwen_avg_values = calculate_values_avg(qwen_posts)

    print("GPT averages:", gpt_avg_values)
    print("\nMistral averages:", mistral_avg_values)
    print("\nDeepSeek averages:", deepseek_avg_values)
    print("\nQwen averages:", qwen_avg_values)

    # Spider plot individuali
    plot_spider_chart(
        avg_values={"GPT-4.1-mini": gpt_avg_values},
        output_path="spider_gpt.png",
        title="GPT-4.1-mini - Schwartz Values",
    )

    plot_spider_chart(
        avg_values={"Mistral-14b": mistral_avg_values},
        output_path="spider_mistral.png",
        title="Mistral-14b - Schwartz Values",
    )

    plot_spider_chart(
        avg_values={"DeepSeek": deepseek_avg_values},
        output_path="spider_deepseek.png",
        title="DeepSeek - Schwartz Values",
    )

    plot_spider_chart(
        avg_values={"Qwen3": qwen_avg_values},
        output_path="spider_qwen.png",
        title="Qwen3 - Schwartz Values",
    )

    # Bar chart individuali
    plot_values_comparison(
        avg_values={"GPT-4.1-mini": gpt_avg_values},
        output_path="bar_gpt.png",
    )

    plot_values_comparison(
        avg_values={"Mistral-14b": mistral_avg_values},
        output_path="bar_mistral.png",
    )

    plot_values_comparison(
        avg_values={"DeepSeek": deepseek_avg_values},
        output_path="bar_deepseek.png",
    )

    plot_values_comparison(
        avg_values={"Qwen3": qwen_avg_values},
        output_path="bar_qwen.png",
    )

    # Spider plot combinato
    plot_spider_chart(
        avg_values={
            "GPT-4.1-mini": gpt_avg_values,
            "Mistral-14b": mistral_avg_values,
            "DeepSeek": deepseek_avg_values,
            "Qwen3": qwen_avg_values,
        },
        output_path="spider_4_models.png",
    )

    # Bar chart combinato
    plot_values_comparison(
        avg_values={
            "GPT-4.1-mini": gpt_avg_values,
            "Mistral-14b": mistral_avg_values,
            "DeepSeek": deepseek_avg_values,
            "Qwen3": qwen_avg_values,
        },
        output_path="comparison_4_models.png",
    )


if __name__ == "__main__":
    main()
