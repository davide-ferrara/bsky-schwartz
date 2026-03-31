import matplotlib.pyplot as plt
import seaborn as sns
from matplotlib.patches import Patch
from typing import Optional


SCHWARTZ_CLUSTERS = {
    "Openness to Change": {
        "values": [
            "Independent thoughts",
            "Independent actions",
            "Stimulation",
            "Pleasure",
        ],
        "color": "#FF69B4",
    },
    "Self-Enhancement": {
        "values": ["Achievement", "Power", "Wealth", "Reputation"],
        "color": "#FFA500",
    },
    "Conservation": {
        "values": [
            "Personal security",
            "Societal security",
            "Tradition",
            "Lawfulness",
            "Respect",
            "Humility",
        ],
        "color": "#4169E1",
    },
    "Self-Transcendence": {
        "values": ["Caring", "Responsibility", "Equality", "Nature", "Tolerance"],
        "color": "#228B22",
    },
}

MODEL_COLORS = {
    "GPT": "#27AE60",
    "GPT-4.1-mini": "#27AE60",
    "Mistral": "#FF6B35",
    "Mistral-14b": "#FF6B35",
    "ministral-14b": "#FF6B35",
    "DeepSeek": "#2E86AB",
    "DeepSeek-V3": "#2E86AB",
    "deepseek-v3": "#2E86AB",
    "Qwen": "#8E44AD",
    "Qwen3": "#8E44AD",
    "Qwen3-VL-30B": "#8E44AD",
}


def plot_values_comparison(
    avg_values: dict[str, dict[str, float]],
    output_path: str,
    title: str = None,
    figsize: tuple[int, int] = None,
) -> None:
    """
    Genera un bar chart raggruppato per confrontare valori Schwartz medi.

    Args:
        avg_values: Dizionario {modello: {valore: media}}
                    Es: {"GPT": {"Reputation": 2.0, ...}, "Gemini": {...}}
        output_path: Path dove salvare il PNG
        title: Titolo del grafico (se None, genera automaticamente)
        figsize: Dimensioni figura (se None, calcola in base a N modelli)
    """
    models = list(avg_values.keys())
    n_models = len(models)

    # Auto-genera titolo
    if title is None:
        if len(models) <= 3:
            title = " vs ".join(models)
        else:
            title = (
                f"{models[0]} vs {models[1]} vs {models[2]} (+{len(models) - 3} more)"
            )

    # Auto-calcola dimensioni figura
    if figsize is None:
        height = 10 + (n_models * 0.8)
        figsize = (14, height)

    sns.set_theme(style="whitegrid")

    fig, ax = plt.subplots(figsize=figsize)

    bar_height = 0.35
    group_spacing = 0.2
    cluster_gap = 1.2

    current_y = 0
    y_ticks = []
    y_labels = []

    cluster_info_to_plot = []

    clusters_to_process = list(SCHWARTZ_CLUSTERS.items())[::-1]

    for cluster_name, cluster_data in clusters_to_process:
        available_vals = [
            v for v in cluster_data["values"] if v in avg_values[models[0]]
        ]
        if not available_vals:
            continue

        cluster_start_y = current_y

        for val_name in available_vals[::-1]:
            for i, model in enumerate(models):
                val_score = avg_values[model].get(val_name, 0)

                # Try exact match first, then capitalize first letter
                model_color = MODEL_COLORS.get(model)
                if model_color is None:
                    # Try with first part capitalized
                    first_part = model.split("-")[0]
                    model_color = MODEL_COLORS.get(first_part.capitalize())
                if model_color is None:
                    model_color = f"C{i}"

                pos = current_y + (i * bar_height)

                ax.barh(
                    pos,
                    val_score,
                    height=bar_height,
                    color=model_color,
                    edgecolor="white",
                    alpha=0.9,
                    label=model if current_y == 0 else "",
                )

                ax.text(
                    val_score + 0.05,
                    pos,
                    f"{val_score:.1f}",
                    va="center",
                    ha="left",
                    fontsize=9,
                    fontweight="bold",
                )

            y_ticks.append(current_y + (bar_height * (n_models - 1) / 2))
            y_labels.append(val_name)

            current_y += (n_models * bar_height) + group_spacing

        cluster_end_y = current_y - group_spacing
        cluster_info_to_plot.append(
            {
                "name": cluster_name,
                "center": (cluster_start_y + cluster_end_y - bar_height) / 2,
                "color": cluster_data["color"],
                "start": cluster_start_y,
                "end": cluster_end_y,
            }
        )

        current_y += cluster_gap

    # Background colorato per cluster
    for info in cluster_info_to_plot:
        start_y = info["start"] - bar_height
        end_y = info["end"] + bar_height
        ax.axhspan(start_y, end_y, color=info["color"], alpha=0.08, zorder=0)

    ax.set_yticks(y_ticks)
    ax.set_yticklabels(y_labels, fontweight="bold", fontsize=10, color="black")
    ax.set_xlabel("Average Value (0-6)", fontsize=11, fontweight="bold")
    ax.set_title(title, fontsize=16, fontweight="bold", pad=25)
    ax.set_xlim(0, 6)

    # Legenda modelli in alto a destra
    handles, labels = ax.get_legend_handles_labels()
    by_label = dict(zip(labels, handles))
    ax.legend(
        by_label.values(), by_label.keys(), loc="upper right", frameon=True, shadow=True
    )

    # Legenda cluster in alto a sinistra
    cluster_legend_elements = [
        Patch(facecolor="#228B22", alpha=0.3, label="Self-Transcendence"),
        Patch(facecolor="#4169E1", alpha=0.3, label="Conservation"),
        Patch(facecolor="#FFA500", alpha=0.3, label="Self-Enhancement"),
        Patch(facecolor="#FF69B4", alpha=0.3, label="Openness to Change"),
    ]
    fig.legend(
        handles=cluster_legend_elements,
        loc="upper left",
        frameon=True,
        shadow=True,
        fontsize=9,
    )

    plt.tight_layout()
    plt.savefig(output_path, dpi=150, bbox_inches="tight")
    plt.close()

    print(f"Chart saved to: {output_path}")
