import matplotlib.pyplot as plt
import os

# Пути к файлам с данными
reader_time_path = "/sysflow/sf-aux/check_perf/time_reader"
aggregator_time_path = "/sysflow/sf-aux/check_perf/time_aggregator"
producer_time_path = "/sysflow/sf-aux/check_perf/time_producer"

# Функция для чтения данных из лог-файла
def read_time_data(file_path):
    if not os.path.exists(file_path):
        print(f"File not found: {file_path}")
        return []

    with open(file_path, "r") as file:
        data = []
        for line in file:
            try:
                # Каждая строка должна содержать числовое значение времени выполнения (мкс)
                data.append(float(line.strip()))
            except ValueError:
                print(f"Skipping invalid line: {line.strip()}")
        return data

# Чтение данных из логов
reader_times = read_time_data(reader_time_path)
aggregator_times = read_time_data(aggregator_time_path)
producer_times = read_time_data(producer_time_path)

# Проверка, есть ли данные для графика
if not reader_times and not aggregator_times and not producer_times:
    print("No data to plot. Exiting.")
    exit()

# Генерация временной шкалы
reader_time_axis = [sum(reader_times[:i]) / 1_000_000 for i in range(1, len(reader_times) + 1)]
aggregator_time_axis = [sum(aggregator_times[:i]) / 1_000_000 for i in range(1, len(aggregator_times) + 1)]
producer_time_axis = [sum(producer_times[:i]) / 1_000_000 for i in range(1, len(producer_times) + 1)]

# Преобразование времени выполнения в микросекунды
reader_times_us = reader_times
aggregator_times_us = aggregator_times
producer_times_us = producer_times

# Построение графиков
fig, axs = plt.subplots(3, 1, figsize=(12, 10), sharex=True)

# График для времени выполнения reader
if reader_times:
    axs[0].plot(reader_time_axis, reader_times_us, label="Reader Goroutine", color="blue")
    axs[0].set_title("Reader Goroutine")
    axs[0].set_ylabel("Cycle Time (µs)")
    axs[0].grid(True)

# График для времени выполнения aggregator
if aggregator_times:
    axs[1].plot(aggregator_time_axis, aggregator_times_us, label="Aggregator Goroutine", color="green")
    axs[1].set_title("Aggregator Goroutine")
    axs[1].set_ylabel("Cycle Time (µs)")
    axs[1].grid(True)

# График для времени выполнения producer
if producer_times:
    axs[2].plot(producer_time_axis, producer_times_us, label="Producer Goroutine", color="red")
    axs[2].set_title("Producer Goroutine")
    axs[2].set_xlabel("Total Elapsed Time (seconds)")
    axs[2].set_ylabel("Cycle Time (µs)")
    axs[2].grid(True)

# Настройка графиков
plt.tight_layout()
plt.savefig("goroutine_cycle_times_stacked.png")
print("Goroutine cycle times graph saved as goroutine_cycle_times_stacked.png")
