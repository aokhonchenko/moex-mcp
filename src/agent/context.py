#!/usr/bin/env python3
"""
Модуль управления контекстом сессии.

Читает файлы точечно через partial_reader, экономя токены модели.
Вместо чтения целого файла — только нужные секции, заголовки, первые N строк.

Создан: сессия 22 (2026-07-06)
Цель: оптимизировать работу агента с моделью через точечное воздействие.

Использование:
    from context import SessionContext

    ctx = SessionContext()
    ctx.load()  # загружает минимальный контекст
    ctx.get_state()  # возвращает словарь состояния
"""

import sys
import os
from pathlib import Path

# Добавляем путь к tools для импорта partial_reader
sys.path.insert(0, str(Path(__file__).parent.parent / 'tools'))

try:
    from partial_reader import read_head, read_headers, read_section, read_summary
except ImportError:
    # Fallback: читаем файл напрямую если partial_reader недоступен
    def read_head(filepath, n=30):
        with open(filepath, 'r', encoding='utf-8') as f:
            return ''.join(f.readline() for _ in range(n))
    
    def read_headers(filepath):
        with open(filepath, 'r', encoding='utf-8') as f:
            return '\n'.join(line.rstrip() for line in f if line.startswith('#'))
    
    def read_section(filepath, section_name):
        # Упрощённая версия
        with open(filepath, 'r', encoding='utf-8') as f:
            return f.read()
    
    def read_summary(filepath, context_lines=2):
        return read_head(filepath, 50)


class SessionContext:
    """
    Управляет контекстом сессии.
    
    Загружает только необходимую информацию, используя точечное чтение.
    Каждый файл читается оптимальным способом:
    - Маленькие (≤40 строк) — целиком
    - Средние (40-100 строк) — заголовки + нужные секции
    - Большие (>100 строк) — только заголовки + summary
    """
    
    def __init__(self, root=None):
        """
        Args:
            root: корневая директория проекта (по умолчанию — текущая)
        """
        self.root = Path(root) if root else Path.cwd()
        self._state = {}
        self._loaded = False
    
    def load(self):
        """Загружает минимальный контекст сессии."""
        self._state = {
            'global_target': self._read_optimized('GLOBAL_TARGET.md'),
            'last_session': self._read_optimized('state/last_session.md'),
            'current_plan': self._read_optimized('state/current_plan.md'),
            'external_messages': self._read_optimized('state/external_messages.md'),
            'active_tasks': self._read_optimized('tasks/active.md'),
        }
        self._loaded = True
        return self
    
    def _read_optimized(self, rel_path):
        """
        Читает файл оптимальным способом, основываясь на размере.
        
        Returns:
            str: содержимое файла (полное или частичное)
        """
        filepath = self.root / rel_path
        
        if not filepath.exists():
            return f"[Файл не найден: {rel_path}]"
        
        try:
            # Определяем размер файла
            size = filepath.stat().st_size
            with open(filepath, 'r', encoding='utf-8') as f:
                line_count = sum(1 for _ in f)
            
            # Стратегия чтения по размеру
            if line_count <= 40:
                # Маленький файл — читаем целиком
                with open(filepath, 'r', encoding='utf-8') as f:
                    return f.read()
            elif line_count <= 100:
                # Средний файл — заголовки + первые строки
                return read_summary(str(filepath), context_lines=2)
            else:
                # Большой файл — только заголовки
                return read_headers(str(filepath))
                
        except Exception as e:
            return f"[Ошибка чтения {rel_path}: {e}]"
    
    def get_state(self):
        """Возвращает загруженное состояние."""
        if not self._loaded:
            self.load()
        return self._state
    
    def get_file_summary(self, rel_path):
        """
        Возвращает краткую сводку файла.
        
        Args:
            rel_path: относительный путь к файлу
            
        Returns:
            str: заголовки + первые строки секций
        """
        filepath = self.root / rel_path
        if not filepath.exists():
            return f"[Файл не найден: {rel_path}]"
        return read_summary(str(filepath), context_lines=2)
    
    def get_file_headers(self, rel_path):
        """
        Возвращает только заголовки файла.
        
        Args:
            rel_path: относительный путь к файлу
            
        Returns:
            str: заголовки markdown (# и ##)
        """
        filepath = self.root / rel_path
        if not filepath.exists():
            return f"[Файл не найден: {rel_path}]"
        return read_headers(str(filepath))
    
    def get_section(self, rel_path, section_name):
        """
        Возвращает特定ную секцию файла.
        
        Args:
            rel_path: относительный путь к файлу
            section_name: имя секции (подстрока заголовка)
            
        Returns:
            str: содержимое секции
        """
        filepath = self.root / rel_path
        if not filepath.exists():
            return f"[Файл не найден: {rel_path}]"
        return read_section(str(filepath), section_name)
    
    def get_file_info(self, rel_path):
        """
        Возвращает информацию о файле.
        
        Args:
            rel_path: относительный путь к файлу
            
        Returns:
            dict: {path, size, lines, strategy}
        """
        filepath = self.root / rel_path
        if not filepath.exists():
            return {'path': rel_path, 'exists': False}
        
        size = filepath.stat().st_size
        with open(filepath, 'r', encoding='utf-8') as f:
            line_count = sum(1 for _ in f)
        
        # Определяем стратегию чтения
        if line_count <= 40:
            strategy = 'full'
        elif line_count <= 100:
            strategy = 'summary'
        else:
            strategy = 'headers_only'
        
        return {
            'path': rel_path,
            'exists': True,
            'size': size,
            'lines': line_count,
            'strategy': strategy
        }


def print_context_stats(ctx):
    """Выводит статистику по загруженному контексту."""
    state = ctx.get_state()
    
    print("=== Статистика контекста сессии ===\n")
    
    total_chars = 0
    for key, content in state.items():
        chars = len(content)
        total_chars += chars
        lines = content.count('\n') + 1
        print(f"  {key}: {lines} строк, {chars} символов")
    
    print(f"\n  Итого: {total_chars} символов")
    print(f"  Примерно: {total_chars // 4} токенов (оценка)")


def main():
    """CLI для тестирования модуля."""
    if len(sys.argv) < 2:
        print("Использование: python context.py <команда> [параметры]")
        print()
        print("Команды:")
        print("  load          — загрузить и показать контекст")
        print("  stats         — статистика по контексту")
        print("  summary FILE  — сводка файла")
        print("  headers FILE  — заголовки файла")
        print("  section FILE X — секция файла")
        print("  info FILE     — информация о файле")
        sys.exit(1)
    
    cmd = sys.argv[1]
    ctx = SessionContext()
    
    if cmd == 'load':
        ctx.load()
        state = ctx.get_state()
        for key, content in state.items():
            print(f"\n{'='*60}")
            print(f"  {key}")
            print(f"{'='*60}")
            print(content[:500])  # Первые 500 символов
            if len(content) > 500:
                print(f"\n  ... ({len(content)} всего символов)")
    
    elif cmd == 'stats':
        ctx.load()
        print_context_stats(ctx)
    
    elif cmd == 'summary':
        if len(sys.argv) < 3:
            print("Ошибка: укажите файл")
            sys.exit(1)
        print(ctx.get_file_summary(sys.argv[2]))
    
    elif cmd == 'headers':
        if len(sys.argv) < 3:
            print("Ошибка: укажите файл")
            sys.exit(1)
        print(ctx.get_file_headers(sys.argv[2]))
    
    elif cmd == 'section':
        if len(sys.argv) < 4:
            print("Ошибка: укажите файл и имя секции")
            sys.exit(1)
        print(ctx.get_section(sys.argv[2], sys.argv[3]))
    
    elif cmd == 'info':
        if len(sys.argv) < 3:
            print("Ошибка: укажите файл")
            sys.exit(1)
        info = ctx.get_file_info(sys.argv[2])
        for k, v in info.items():
            print(f"  {k}: {v}")
    
    else:
        print(f"Неизвестная команда: {cmd}")
        sys.exit(1)


if __name__ == '__main__':
    main()
