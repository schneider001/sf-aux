import os

def calculate_average_time(file_path):
    """
    Считывает значения времени из файла, рассчитывает и возвращает среднее время.
    :param file_path: Путь к файлу с данными времени.
    :return: Среднее время выполнения в микросекундах или None, если файл пуст.
    """
    try:
        if not os.path.exists(file_path):
            print(f"Файл {file_path} не найден.")
            return None

        with open(file_path, 'r') as file:
            times = [int(line.strip()) for line in file if line.strip().isdigit()]

        if not times:
            print(f"Файл {file_path} пуст.")
            return None

        return sum(times) / len(times)
    except Exception as e:
        print(f"Ошибка при обработке файла {file_path}: {e}")
        return None


def main():
    files = [
        "/sysflow/sf-aux/check_perf/time_reader",
        "/sysflow/sf-aux/check_perf/time_aggregator",
        "/sysflow/sf-aux/check_perf/time_producer"
    ]

    for file_path in files:
        avg_time = calculate_average_time(file_path)
        if avg_time is not None:
            print(f"Среднее время выполнения из {file_path}: {avg_time:.2f} микросекунд")


if __name__ == "__main__":
    main()
