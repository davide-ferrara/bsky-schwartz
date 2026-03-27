import json
import random


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


def generate_mock_data(
    input_file: str, output_file: str, model: str = "qwen/qwen-2.5-72b-instruct"
):
    with open(input_file) as f:
        data = json.load(f)

    for post in data:
        values = {k: random.randint(0, 6) for k in VALUES}
        post["values"] = values
        post["values_arr"] = [values[k] for k in VALUES]
        post["score"] = sum(post["values_arr"])
        post["model"] = model
        post["stats"] = {
            "response_time_ms": random.randint(5000, 30000),
            "tokens_used": random.randint(1000, 2000),
            "cost_usd": round(random.uniform(0.0001, 0.003), 7),
        }

    with open(output_file, "w") as f:
        json.dump(data, f, indent=2)

    print(f"Generated {output_file} with {len(data)} posts from model {model}")


if __name__ == "__main__":
    generate_mock_data("post_data_gpt.json", "post_data_qwen.json")
