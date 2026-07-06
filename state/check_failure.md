# Диагностика упавших проверок

- Попытка исправления: 1/2
- Команда проверки: `c:\_dev\own\pet\.venv\Scripts\python.exe -m pytest`

Проверки упали. Это не финальный результат сессии: исправь ошибки в текущем временном worktree и снова обнови обязательные файлы сессии.

## Вывод проверок

```text
... 21 earlier output lines omitted; tail follows ...
tests\test_sleep_memory.py ...                                           [ 99%]
tests\test_validation_repairs.py ..                                      [100%]

================================== FAILURES ===================================
____________________ TestReplaceRegex.test_regex_multiline ____________________

self = <test_apply_patch.TestReplaceRegex object at 0x000002277C4BB550>

    def test_regex_multiline(self):
        content = 'a=1\nb=2\na=3\n'
        path = _make_temp_file(content)
        try:
            result = replace_regex(path, r'^a=\d+', 'a=0')
            assert result.applied is True
>           assert result.changes == 2
E           AssertionError: assert 1 == 2
E            +  where 1 = PatchResult(path='C:\\Users\\ohotNik\\AppData\\Local\\Temp\\tmpqtifxpk8.py', applied=True, operation='regex', changes=1, preview='�������� 1 ��������� �� ������� ^a=\\d+', error=None).changes

tests\test_apply_patch.py:148: AssertionError
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
FAILED tests/test_apply_patch.py::TestReplaceRegex::test_regex_multiline - As...
======================== 1 failed, 204 passed in 2.46s ========================
```
