import io
import pandas as pd
import matplotlib.pyplot as plt

df = pd.read_csv('plot.csv',  usecols=['Time', 'active-actors', 'submitted', 'completed'])
#df["Time"] = pd.to_datetime(df["Time"])

fig, ax = plt.subplots(figsize=(12,5))
ax.set_title('actor benchmark')
ax2 = ax.twinx()
df.plot(x='Date', y='Task', kind='line', ax=ax, alpha=0.2, color='black')
ax.set_xticklabels(df['Date'].dt.date, fontsize=8)
ax.get_legend().remove()
ax2.plot(ax.get_xticks(), df['active-actors'], color='black')
ax.set_ylabel('task')
ax.set_yticks([0, 3000])
ax2.set_ylabel('actors')
plt.tight_layout()
plt.show()


plt.show()