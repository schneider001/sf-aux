# import matplotlib.pyplot as plt
# import os

# # Пути к файлам с логами
# records_fill_path = "/sysflow/sf-aux/check_perf/records_fill"
# events_fill_path = "/sysflow/sf-aux/check_perf/events_fill"

# # Функция для чтения данных из лог-файла
# def read_fill_data(file_path):
#     if not os.path.exists(file_path):
#         print(f"File not found: {file_path}")
#         return []

#     with open(file_path, "r") as file:
#         data = []
#         for line in file:
#             try:
#                 # Каждая строка должна содержать числовое значение заполненности
#                 data.append(int(line.strip()))
#             except ValueError:
#                 print(f"Skipping invalid line: {line.strip()}")
#         return data

# # Чтение данных из логов
# records_fill = read_fill_data(records_fill_path)
# events_fill = read_fill_data(events_fill_path)

# # Проверка, есть ли данные для графиков
# if not records_fill and not events_fill:
#     print("No data to plot. Exiting.")
#     exit()

# # Генерация временной шкалы
# records_time = [i * 0.01 for i in range(len(records_fill))]  # Интервал 0.01 секунды
# events_time = [i * 0.01 for i in range(len(events_fill))]

# # Построение графика для records_fill
# if records_fill:
#     plt.figure(figsize=(12, 6))
#     plt.plot(records_time, records_fill, label="Records Fill", color="blue")
#     plt.title("Records Fill Levels Over Time")
#     plt.xlabel("Time (seconds)")
#     plt.ylabel("Channel Fill Level")
#     plt.legend()
#     plt.grid(True)
#     plt.tight_layout()
#     plt.savefig("records_fill_graph.png")
#     print("Records fill graph saved as records_fill_graph.png")

# # Построение графика для events_fill
# if events_fill:
#     plt.figure(figsize=(12, 6))
#     plt.plot(events_time, events_fill, label="Events Fill", color="orange")
#     plt.title("Events Fill Levels Over Time")
#     plt.xlabel("Time (seconds)")
#     plt.ylabel("Channel Fill Level")
#     plt.legend()
#     plt.grid(True)
#     plt.tight_layout()
#     plt.savefig("events_fill_graph.png")
#     print("Events fill graph saved as events_fill_graph.png")

import matplotlib.pyplot as plt
import os

# Пути к файлам с логами
records_fill_path = "/sysflow/sf-aux/check_perf/records_fill"
events_fill_path = "/sysflow/sf-aux/check_perf/events_fill"

# Функция для чтения данных из лог-файла
def read_fill_data(file_path):
    if not os.path.exists(file_path):
        print(f"File not found: {file_path}")
        return []

    with open(file_path, "r") as file:
        data = []
        for line in file:
            try:
                # Каждая строка должна содержать числовое значение заполненности
                data.append(int(line.strip()))
            except ValueError:
                print(f"Skipping invalid line: {line.strip()}")
        return data

# Чтение данных из логов
records_fill = read_fill_data(records_fill_path)
events_fill = read_fill_data(events_fill_path)

# Проверка, есть ли данные для графиков
if not records_fill and not events_fill:
    print("No data to plot. Exiting.")
    exit()

# Генерация временной шкалы
records_time = [i * 0.01 for i in range(len(records_fill))]  # Интервал 0.01 секунды
events_time = [i * 0.01 for i in range(len(events_fill))]

# Создание фигуры с 2 подграфиками (2 строки, 1 столбец)
fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(12, 12))

# Построение графика для records_fill
if records_fill:
    ax1.plot(records_time, records_fill, label="Records Fill", color="blue")
    ax1.set_title("Records Fill Levels Over Time")
    ax1.set_xlabel("Time (seconds)")
    ax1.set_ylabel("Channel Fill Level")
    ax1.legend()
    ax1.grid(True)

# Построение графика для events_fill
if events_fill:
    ax2.plot(events_time, events_fill, label="Events Fill", color="orange")
    ax2.set_title("Events Fill Levels Over Time")
    ax2.set_xlabel("Time (seconds)")
    ax2.set_ylabel("Channel Fill Level")
    ax2.legend()
    ax2.grid(True)

# Настройка макета для улучшенного отображения
plt.tight_layout()

# Сохранение итогового графика
plt.savefig("combined_fill_graph.png")
print("Combined fill graph saved as combined_fill_graph.png")
