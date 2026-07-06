from scripts import run_snapshots


def test_preserve_session_snapshot_copies_worktree_and_prunes_old_snapshots(tmp_path):
    runs_dir = tmp_path / "runs"
    source = tmp_path / "worktree"
    source.mkdir()
    (source / "state").mkdir()
    (source / "state" / "last_session.md").write_text("готово\n", encoding="utf-8")
    (source / ".git").write_text("gitdir: ignored\n", encoding="utf-8")
    snapshots = runs_dir / "snapshots"
    old_one = snapshots / "session-0001-applied"
    old_two = snapshots / "session-0002-failed"
    old_one.mkdir(parents=True)
    old_two.mkdir()
    (old_one / "old.md").write_text("1", encoding="utf-8")
    (old_two / "old.md").write_text("2", encoding="utf-8")

    target = run_snapshots.preserve_session_snapshot(source, runs_dir, "0003", applied=False)

    assert target == snapshots / "session-0003-failed"
    assert (target / "state" / "last_session.md").read_text(encoding="utf-8") == "готово\n"
    assert not (target / ".git").exists()
    assert "завершилась ошибкой" in (target / "RUN_SNAPSHOT.md").read_text(encoding="utf-8")
    assert not old_one.exists()
    assert old_two.exists()


def test_prune_snapshots_keeps_directory_when_keep_is_zero(tmp_path):
    snapshots = tmp_path / "snapshots"
    old = snapshots / "session-0001-applied"
    old.mkdir(parents=True)

    run_snapshots.prune_snapshots(snapshots, keep=0)

    assert old.exists()