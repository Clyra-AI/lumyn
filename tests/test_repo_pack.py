from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[1]


class RepoPackTests(unittest.TestCase):
    def test_operating_pack_exists(self) -> None:
        for relative_path in [
            "AGENTS.md",
            "WORKFLOW.md",
            "README.md",
            "docs/product/prd.md",
            ".factory/artifacts/prd-to-plan/lumyn-mvp/context-brief.json",
            ".factory/artifacts/prd-to-plan/lumyn-mvp/execution-plan.json",
        ]:
            self.assertTrue((ROOT / relative_path).exists(), relative_path)

    def test_prd_references_are_repo_relative(self) -> None:
        prd = (ROOT / "docs/product/prd.md").read_text()
        self.assertIn("Lumyn OSS MVP", prd)
        self.assertNotIn("/" + "Users/", prd)


if __name__ == "__main__":
    unittest.main()
