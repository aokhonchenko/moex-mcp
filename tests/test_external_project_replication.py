from pathlib import Path

from scripts import external_projects, session_transaction
from test_session_transaction import FakeRunner, make_root


def test_external_project_dirs_include_nested_repos_and_ignored_dirs(tmp_path):
    root = make_root(tmp_path)
    nested = root / "projects" / "external-repo"
    ignored = root / "projects" / "foundation-finance"
    tracked = root / "projects" / "notes"
    (nested / ".git").mkdir(parents=True)
    ignored.mkdir(parents=True)
    tracked.mkdir(parents=True)

    def runner(args, cwd):
        if list(args)[:3] == ["git", "check-ignore", "-q"]:
            return session_transaction.CommandResult(0 if args[3] == "projects/foundation-finance" else 1, "", "")
        return session_transaction.CommandResult(0, "", "")

    result = external_projects.external_project_dirs(root, runner)

    assert [path.relative_to(root).as_posix() for path in result] == [
        "projects/external-repo",
        "projects/foundation-finance",
    ]


def test_replicate_external_projects_copies_ignored_project_without_tracking_it(tmp_path):
    main_parent = tmp_path / "main"
    worktree_parent = tmp_path / "worktree"
    main_parent.mkdir()
    worktree_parent.mkdir()
    root = make_root(main_parent)
    worktree = make_root(worktree_parent)
    project = worktree / "projects" / "foundation-finance"
    project.mkdir(parents=True)
    (project / "README.md").write_text("dashboard\n", encoding="utf-8")
    (project / ".git").mkdir()
    (project / "node_modules").mkdir()
    (project / "node_modules" / "large.txt").write_text("skip\n", encoding="utf-8")

    def runner(args, cwd):
        if list(args)[:3] == ["git", "check-ignore", "-q"]:
            return session_transaction.CommandResult(0 if args[3] == "projects/foundation-finance" else 1, "", "")
        return session_transaction.CommandResult(0, "", "")

    replicated = external_projects.replicate_external_projects(worktree, root, runner)

    assert replicated == [Path("projects/foundation-finance")]
    assert (root / "projects" / "foundation-finance" / "README.md").read_text(encoding="utf-8") == "dashboard\n"
    assert (root / "projects" / "foundation-finance" / ".git").exists()
    assert not (root / "projects" / "foundation-finance" / "node_modules").exists()


def test_run_transaction_replicates_external_projects_after_success(tmp_path):
    root = make_root(tmp_path)
    runs_dir = tmp_path / "runs"
    runner = FakeRunner(root, create_external_project=True)

    session_transaction.run_transaction(
        root=root,
        agent_command="agent --ok",
        runs_dir=runs_dir,
        runner=runner,
    )

    assert (root / "projects" / "foundation-finance" / "README.md").read_text(encoding="utf-8") == "created\n"
    assert (root / "projects" / "foundation-finance" / ".git").exists()
