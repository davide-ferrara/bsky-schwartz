import json
from charts import plot_values_comparison


def load_json(file_path: str):
    if not file_path.endswith("json"):
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
    # Carica i dati
    gpt_posts = load_json("../post_gpt-4.1-mini_20260330232334.json")
    gemini_posts = load_json(
        "../post_gemini-3.1-flash-lite-preview_20260330232401.json"
    )
    mistral_posts = load_json("../post_mistral-small-2603_20260330234657.json")

    if not gpt_posts or not gemini_posts or not mistral_posts:
        print("Error loading data")
        return

    print("Data loaded successfully")

    # Calcola le medie per ogni modello
    gpt_avg_values = calculate_values_avg(gpt_posts)
    gemini_avg_values = calculate_values_avg(gemini_posts)
    mistral_avg_values = calculate_values_avg(mistral_posts)

    print("GPT averages:", gpt_avg_values)
    print("\nGemini averages:", gemini_avg_values)
    print("\nMistral averages:", mistral_avg_values)

    # Genera grafico con 3 modelli
    plot_values_comparison(
        avg_values={
            "GPT-4.1-mini": gpt_avg_values,
            "Gemini-3.1-flash": gemini_avg_values,
            "Mistral-Small": mistral_avg_values,
        },
        output_path="comparison_3models.png",
    )


if __name__ == "__main__":
    main()
