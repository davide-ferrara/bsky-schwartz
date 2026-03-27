import json

from charts import plot_comparison_chart, plot_radar_chart


VALUES = [
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


def load_json(file_path: str):
    if not file_path.endswith("json"):
        return None

    with open(file_path) as f:
        data = json.load(f)

    return data


def compute_averages(posts: list) -> dict:
    values_avg = {v: 0.0 for v in VALUES}
    for post in posts:
        for value in VALUES:
            values_avg[value] += float(post["values"][value])
    for value in VALUES:
        values_avg[value] /= len(posts)
    return values_avg


def main():
    gpt_data = load_json("post_data_gpt.json")
    qwen_data = load_json("post_data_qwen.json")

    if gpt_data is None or qwen_data is None:
        print("Error loading data")
        return

    print("Data loaded successfully")

    res = {
        gpt_data[0]["model"]: compute_averages(gpt_data),
        qwen_data[0]["model"]: compute_averages(qwen_data),
    }

    print("Averages:")
    for model, values in res.items():
        print(f"  {model}:")
        for k, v in values.items():
            print(f"    {k}: {v:.2f}")

    plot_comparison_chart(res, "comparison.png")
    plot_radar_chart(res, "radar.png")


if __name__ == "__main__":
    main()
