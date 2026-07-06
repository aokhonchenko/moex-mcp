# Диагностика упавших проверок

- Попытка исправления: 1/2
- Команда проверки: `c:\_dev\own\pet\.venv\Scripts\python.exe -m pytest`

Проверки упали. Это не финальный результат сессии: исправь ошибки в текущем временном worktree и снова обнови обязательные файлы сессии.

## Вывод проверок

```text
... 98 earlier output lines omitted; tail follows ...
    def test_method_in_class(self):
        content = textwrap.dedent('''\
            class MyClass:
                def method(self, x):
                    """�����."""
                    return x
        ''')
        path = _make_temp_file(content)
        try:
            result = read_func(path, 'method')
>           assert 'def method' in result.content
E           assert 'def method' in '        """�����."""\n        return x\n'
E            +  where '        """�����."""\n        return x\n' = ReadResult(path='C:\\Users\\ohotNik\\AppData\\Local\\Temp\\tmpvf9vavin.py', content='        """�����."""\n        return x\n', lines_read=2, method='func[method]', truncated=False).content

tests\test_reader.py:244: AssertionError
=============================== tests coverage ================================
______________ coverage: platform win32, python 3.11.15-final-0 _______________

Name                             Stmts   Miss Branch BrPart  Cover   Missing
----------------------------------------------------------------------------
scripts\command_runners.py          68      2      8      1    96%   92-93, 104->107
scripts\file_tools.py               42      2     14      2    93%   54, 74
scripts\llm_client.py               62      5     10      1    92%   26, 74-75, 79-80
scripts\run_agent.py               139      6     34      5    94%   14, 136-137, 146, 219, 221->224, 232
scripts\run_session.py              91      2     22      2    96%   214, 253
scripts\run_snapshots.py            35      3      6      1    90%   60, 69-70
scripts\session_transaction.py     281     12     74     10    94%   14, 97, 102, 114, 124, 246->exit, 274->263, 344, 350, 390-391, 411-412, 481
scripts\sleep_memory.py             84     18     18      3    79%   22, 51-53, 133-136, 140-149, 153
----------------------------------------------------------------------------
TOTAL                              842     50    190     25    93%

2 files skipped due to complete coverage.
Required test coverage of 90% reached. Total coverage: 92.73%
=========================== short test summary info ===========================
FAILED tests/test_reader.py::TestReadLines::test_basic_range - AttributeError...
FAILED tests/test_reader.py::TestReadFunc::test_simple_function - assert 'def...
FAILED tests/test_reader.py::TestReadFunc::test_function_with_decorator - ass...
FAILED tests/test_reader.py::TestReadFunc::test_async_function - assert 'asyn...
FAILED tests/test_reader.py::TestReadFunc::test_method_in_class - assert 'def...
======================== 5 failed, 175 passed in 2.53s ========================
```
