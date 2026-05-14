window.BENCHMARK_DATA = {
  "lastUpdate": 1778763515518,
  "repoUrl": "https://github.com/bmf-san/ggc",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "email": "bmf.infomation@gmail.com",
            "name": "Kenta Takeuchi",
            "username": "bmf-san"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f866282184fa3c67c4981afb05c0c3a82d42fd71",
          "message": "fix(lint): inline reflect.Ptr to satisfy govet (#426)\n\ngolangci-lint's govet check now flags 'inline: Constant reflect.Ptr should be inlined'. Replace the deprecated alias reflect.Ptr with the canonical reflect.Pointer in internal/config/path.go to unblock CI.",
          "timestamp": "2026-05-13T23:20:19+09:00",
          "tree_id": "6eb64abd68228529baf075ce49e1ca69cb5af5cf",
          "url": "https://github.com/bmf-san/ggc/commit/f866282184fa3c67c4981afb05c0c3a82d42fd71"
        },
        "date": 1778682228799,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 32351920,
            "unit": "ns/op\t 1389961 B/op\t    6916 allocs/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32351920,
            "unit": "ns/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1389961,
            "unit": "B/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6916,
            "unit": "allocs/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 32042597,
            "unit": "ns/op\t 1390242 B/op\t    6919 allocs/op",
            "extra": "49 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32042597,
            "unit": "ns/op",
            "extra": "49 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1390242,
            "unit": "B/op",
            "extra": "49 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6919,
            "unit": "allocs/op",
            "extra": "49 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 32019466,
            "unit": "ns/op\t 1388030 B/op\t    6917 allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32019466,
            "unit": "ns/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1388030,
            "unit": "B/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6917,
            "unit": "allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 32076946,
            "unit": "ns/op\t 1394494 B/op\t    6918 allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32076946,
            "unit": "ns/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1394494,
            "unit": "B/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6918,
            "unit": "allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 31961019,
            "unit": "ns/op\t 1395394 B/op\t    6921 allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 31961019,
            "unit": "ns/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1395394,
            "unit": "B/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6921,
            "unit": "allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 32051181,
            "unit": "ns/op\t 1391786 B/op\t    6918 allocs/op",
            "extra": "45 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32051181,
            "unit": "ns/op",
            "extra": "45 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1391786,
            "unit": "B/op",
            "extra": "45 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6918,
            "unit": "allocs/op",
            "extra": "45 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 32418079,
            "unit": "ns/op\t 1474782 B/op\t    7973 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32418079,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474782,
            "unit": "B/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7973,
            "unit": "allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 32141826,
            "unit": "ns/op\t 1473167 B/op\t    7973 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32141826,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1473167,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7973,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 56573850,
            "unit": "ns/op\t 1473026 B/op\t    7972 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 56573850,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1473026,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7972,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 32299386,
            "unit": "ns/op\t 1473312 B/op\t    7972 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32299386,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1473312,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7972,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 32382059,
            "unit": "ns/op\t 1473830 B/op\t    7972 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32382059,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1473830,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7972,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 32047329,
            "unit": "ns/op\t 1473863 B/op\t    7973 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 32047329,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1473863,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7973,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2747,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "428211 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2747,
            "unit": "ns/op",
            "extra": "428211 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "428211 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "428211 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2746,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "411126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2746,
            "unit": "ns/op",
            "extra": "411126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "411126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "411126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2802,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "409700 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2802,
            "unit": "ns/op",
            "extra": "409700 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "409700 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "409700 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2757,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "411268 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2757,
            "unit": "ns/op",
            "extra": "411268 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "411268 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "411268 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2752,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "408471 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2752,
            "unit": "ns/op",
            "extra": "408471 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "408471 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "408471 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2765,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "412894 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2765,
            "unit": "ns/op",
            "extra": "412894 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "412894 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "412894 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 69.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17444660 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 69.44,
            "unit": "ns/op",
            "extra": "17444660 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17444660 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17444660 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 71.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16160823 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 71.52,
            "unit": "ns/op",
            "extra": "16160823 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16160823 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16160823 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 71.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16942453 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 71.23,
            "unit": "ns/op",
            "extra": "16942453 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16942453 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16942453 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 69.85,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16839742 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 69.85,
            "unit": "ns/op",
            "extra": "16839742 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16839742 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16839742 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 71.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17023917 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 71.51,
            "unit": "ns/op",
            "extra": "17023917 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17023917 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17023917 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 71.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16985293 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 71.5,
            "unit": "ns/op",
            "extra": "16985293 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16985293 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16985293 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5831,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "206568 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5831,
            "unit": "ns/op",
            "extra": "206568 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "206568 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "206568 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5821,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "199646 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5821,
            "unit": "ns/op",
            "extra": "199646 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "199646 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "199646 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5818,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "198370 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5818,
            "unit": "ns/op",
            "extra": "198370 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "198370 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "198370 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5832,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "198849 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5832,
            "unit": "ns/op",
            "extra": "198849 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "198849 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "198849 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5914,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "197572 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5914,
            "unit": "ns/op",
            "extra": "197572 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "197572 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "197572 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5841,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "200893 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5841,
            "unit": "ns/op",
            "extra": "200893 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "200893 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "200893 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1447,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "817832 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1447,
            "unit": "ns/op",
            "extra": "817832 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "817832 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "817832 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1425,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "801763 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1425,
            "unit": "ns/op",
            "extra": "801763 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "801763 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "801763 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1429,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "816194 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1429,
            "unit": "ns/op",
            "extra": "816194 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "816194 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "816194 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1453,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "814471 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1453,
            "unit": "ns/op",
            "extra": "814471 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "814471 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "814471 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1431,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "785647 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1431,
            "unit": "ns/op",
            "extra": "785647 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "785647 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "785647 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1427,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "815412 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1427,
            "unit": "ns/op",
            "extra": "815412 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "815412 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "815412 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 943.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1277730 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 943.7,
            "unit": "ns/op",
            "extra": "1277730 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1277730 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1277730 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 947.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1269296 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 947.2,
            "unit": "ns/op",
            "extra": "1269296 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1269296 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1269296 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 916.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1259674 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 916.5,
            "unit": "ns/op",
            "extra": "1259674 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1259674 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1259674 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 912.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1318483 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 912.5,
            "unit": "ns/op",
            "extra": "1318483 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1318483 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1318483 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 912.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1312676 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 912.9,
            "unit": "ns/op",
            "extra": "1312676 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1312676 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1312676 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 915.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1314920 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 915.5,
            "unit": "ns/op",
            "extra": "1314920 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1314920 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1314920 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1013,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1013,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1014,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1014,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1013,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1013,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1013,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1013,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1010,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1010,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1010,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1010,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1489,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "831600 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1489,
            "unit": "ns/op",
            "extra": "831600 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "831600 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "831600 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1491,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "857847 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1491,
            "unit": "ns/op",
            "extra": "857847 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "857847 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "857847 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1501,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "784076 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1501,
            "unit": "ns/op",
            "extra": "784076 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "784076 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "784076 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1498,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "772380 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1498,
            "unit": "ns/op",
            "extra": "772380 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "772380 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "772380 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1478,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "809097 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1478,
            "unit": "ns/op",
            "extra": "809097 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "809097 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "809097 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1476,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "827203 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1476,
            "unit": "ns/op",
            "extra": "827203 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "827203 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "827203 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1043,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1043,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1039,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1039,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1033,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1033,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1037,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1037,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1032,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1032,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1035,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1035,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 870.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1343290 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 870.6,
            "unit": "ns/op",
            "extra": "1343290 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1343290 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1343290 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 863.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1374528 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 863.5,
            "unit": "ns/op",
            "extra": "1374528 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1374528 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1374528 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 863.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1392951 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 863.5,
            "unit": "ns/op",
            "extra": "1392951 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1392951 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1392951 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 866,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1373748 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 866,
            "unit": "ns/op",
            "extra": "1373748 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1373748 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1373748 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 866.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1370161 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 866.8,
            "unit": "ns/op",
            "extra": "1370161 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1370161 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1370161 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 873.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1372224 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 873.2,
            "unit": "ns/op",
            "extra": "1372224 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1372224 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1372224 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1035,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1035,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1039,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1039,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1053,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1053,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1040,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1040,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1046,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1046,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1038,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1038,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.677,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "427370930 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.677,
            "unit": "ns/op",
            "extra": "427370930 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "427370930 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "427370930 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.689,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "412931095 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.689,
            "unit": "ns/op",
            "extra": "412931095 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "412931095 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "412931095 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.729,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "432891896 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.729,
            "unit": "ns/op",
            "extra": "432891896 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "432891896 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "432891896 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.693,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "421358983 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.693,
            "unit": "ns/op",
            "extra": "421358983 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "421358983 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "421358983 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.835,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "416840284 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.835,
            "unit": "ns/op",
            "extra": "416840284 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "416840284 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "416840284 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.686,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "416208638 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.686,
            "unit": "ns/op",
            "extra": "416208638 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "416208638 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "416208638 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1e1831bf292763e61f980cc9bc419d9256f9f9c6",
          "message": "chore(deps): bump golang.org/x/term from 0.42.0 to 0.43.0 (#423)\n\nBumps [golang.org/x/term](https://github.com/golang/term) from 0.42.0 to 0.43.0.\n- [Commits](https://github.com/golang/term/compare/v0.42.0...v0.43.0)\n\n---\nupdated-dependencies:\n- dependency-name: golang.org/x/term\n  dependency-version: 0.43.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-05-14T21:53:44+09:00",
          "tree_id": "a8f24bb983b729e000f6589f1593973e9b9a7ea6",
          "url": "https://github.com/bmf-san/ggc/commit/1e1831bf292763e61f980cc9bc419d9256f9f9c6"
        },
        "date": 1778763515206,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 25107917,
            "unit": "ns/op\t 1389386 B/op\t    6929 allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 25107917,
            "unit": "ns/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1389386,
            "unit": "B/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6929,
            "unit": "allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 25231451,
            "unit": "ns/op\t 1391512 B/op\t    6930 allocs/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 25231451,
            "unit": "ns/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1391512,
            "unit": "B/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6930,
            "unit": "allocs/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 24944240,
            "unit": "ns/op\t 1393139 B/op\t    6930 allocs/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 24944240,
            "unit": "ns/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1393139,
            "unit": "B/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6930,
            "unit": "allocs/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 24889211,
            "unit": "ns/op\t 1394188 B/op\t    6929 allocs/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 24889211,
            "unit": "ns/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1394188,
            "unit": "B/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6929,
            "unit": "allocs/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 25676762,
            "unit": "ns/op\t 1392590 B/op\t    6931 allocs/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 25676762,
            "unit": "ns/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1392590,
            "unit": "B/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6931,
            "unit": "allocs/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 25579182,
            "unit": "ns/op\t 1392129 B/op\t    6929 allocs/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 25579182,
            "unit": "ns/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1392129,
            "unit": "B/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6929,
            "unit": "allocs/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 26081882,
            "unit": "ns/op\t 1474898 B/op\t    7987 allocs/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 26081882,
            "unit": "ns/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474898,
            "unit": "B/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7987,
            "unit": "allocs/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 25903100,
            "unit": "ns/op\t 1474603 B/op\t    7986 allocs/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 25903100,
            "unit": "ns/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474603,
            "unit": "B/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7986,
            "unit": "allocs/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 25566159,
            "unit": "ns/op\t 1474632 B/op\t    7986 allocs/op",
            "extra": "87 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 25566159,
            "unit": "ns/op",
            "extra": "87 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474632,
            "unit": "B/op",
            "extra": "87 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7986,
            "unit": "allocs/op",
            "extra": "87 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 26122862,
            "unit": "ns/op\t 1474930 B/op\t    7986 allocs/op",
            "extra": "86 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 26122862,
            "unit": "ns/op",
            "extra": "86 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474930,
            "unit": "B/op",
            "extra": "86 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7986,
            "unit": "allocs/op",
            "extra": "86 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 25705988,
            "unit": "ns/op\t 1475214 B/op\t    7985 allocs/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 25705988,
            "unit": "ns/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1475214,
            "unit": "B/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7985,
            "unit": "allocs/op",
            "extra": "84 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 26279972,
            "unit": "ns/op\t 1474682 B/op\t    7986 allocs/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 26279972,
            "unit": "ns/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474682,
            "unit": "B/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7986,
            "unit": "allocs/op",
            "extra": "82 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2782,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "437425 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2782,
            "unit": "ns/op",
            "extra": "437425 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "437425 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "437425 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2748,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "401833 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2748,
            "unit": "ns/op",
            "extra": "401833 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "401833 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "401833 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2771,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "411564 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2771,
            "unit": "ns/op",
            "extra": "411564 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "411564 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "411564 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2789,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "399627 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2789,
            "unit": "ns/op",
            "extra": "399627 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "399627 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "399627 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2852,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "392530 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2852,
            "unit": "ns/op",
            "extra": "392530 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "392530 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "392530 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 2831,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "361707 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 2831,
            "unit": "ns/op",
            "extra": "361707 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "361707 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "361707 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 68.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17175375 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 68.39,
            "unit": "ns/op",
            "extra": "17175375 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17175375 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17175375 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 67.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16935516 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 67.9,
            "unit": "ns/op",
            "extra": "16935516 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16935516 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16935516 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 68.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16083283 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 68.14,
            "unit": "ns/op",
            "extra": "16083283 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16083283 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16083283 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 68.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17205057 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 68.38,
            "unit": "ns/op",
            "extra": "17205057 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17205057 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17205057 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 71.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16971920 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 71.59,
            "unit": "ns/op",
            "extra": "16971920 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16971920 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16971920 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17820926 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 73,
            "unit": "ns/op",
            "extra": "17820926 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17820926 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17820926 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5842,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "204729 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5842,
            "unit": "ns/op",
            "extra": "204729 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "204729 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "204729 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5838,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "198489 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5838,
            "unit": "ns/op",
            "extra": "198489 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "198489 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "198489 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5848,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "200116 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5848,
            "unit": "ns/op",
            "extra": "200116 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "200116 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "200116 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5862,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "196936 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5862,
            "unit": "ns/op",
            "extra": "196936 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "196936 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "196936 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5951,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "199597 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5951,
            "unit": "ns/op",
            "extra": "199597 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "199597 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "199597 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 5942,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "197494 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 5942,
            "unit": "ns/op",
            "extra": "197494 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "197494 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "197494 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1436,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "805630 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1436,
            "unit": "ns/op",
            "extra": "805630 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "805630 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "805630 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1435,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "802711 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1435,
            "unit": "ns/op",
            "extra": "802711 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "802711 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "802711 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1428,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "821462 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1428,
            "unit": "ns/op",
            "extra": "821462 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "821462 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "821462 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1431,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "789051 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1431,
            "unit": "ns/op",
            "extra": "789051 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "789051 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "789051 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1428,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "771386 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1428,
            "unit": "ns/op",
            "extra": "771386 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "771386 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "771386 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1437,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "810454 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1437,
            "unit": "ns/op",
            "extra": "810454 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "810454 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "810454 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 923.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1310226 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 923.5,
            "unit": "ns/op",
            "extra": "1310226 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1310226 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1310226 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 913.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1321892 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 913.7,
            "unit": "ns/op",
            "extra": "1321892 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1321892 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1321892 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 918.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1312726 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 918.4,
            "unit": "ns/op",
            "extra": "1312726 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1312726 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1312726 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 920.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1308055 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 920.9,
            "unit": "ns/op",
            "extra": "1308055 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1308055 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1308055 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 913.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1310286 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 913.5,
            "unit": "ns/op",
            "extra": "1310286 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1310286 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1310286 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 916.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1312137 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 916.1,
            "unit": "ns/op",
            "extra": "1312137 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1312137 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1312137 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1010,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1010,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1012,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1012,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1016,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1016,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1015,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1015,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1024,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1024,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1015,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1015,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1480,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "824936 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1480,
            "unit": "ns/op",
            "extra": "824936 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "824936 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "824936 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1473,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "826838 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1473,
            "unit": "ns/op",
            "extra": "826838 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "826838 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "826838 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1487,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "840968 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1487,
            "unit": "ns/op",
            "extra": "840968 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "840968 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "840968 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1481,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "854404 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1481,
            "unit": "ns/op",
            "extra": "854404 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "854404 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "854404 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1454,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "833520 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1454,
            "unit": "ns/op",
            "extra": "833520 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "833520 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "833520 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1450,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "829882 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1450,
            "unit": "ns/op",
            "extra": "829882 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "829882 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "829882 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1043,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1043,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1055,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1055,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1032,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1032,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1028,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1028,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1033,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1033,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1033,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1033,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 874.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1378294 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 874.1,
            "unit": "ns/op",
            "extra": "1378294 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1378294 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1378294 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 857.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1402333 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 857.9,
            "unit": "ns/op",
            "extra": "1402333 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1402333 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1402333 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 854.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1404594 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 854.9,
            "unit": "ns/op",
            "extra": "1404594 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1404594 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1404594 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 873.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1378332 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 873.8,
            "unit": "ns/op",
            "extra": "1378332 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1378332 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1378332 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 876.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1397023 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 876.8,
            "unit": "ns/op",
            "extra": "1397023 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1397023 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1397023 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 855.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1382484 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 855.9,
            "unit": "ns/op",
            "extra": "1382484 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1382484 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1382484 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1042,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1042,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1041,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1041,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1043,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1043,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1040,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1040,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1031,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1031,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1040,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1040,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.732,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "448080499 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.732,
            "unit": "ns/op",
            "extra": "448080499 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "448080499 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "448080499 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.778,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "364040192 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.778,
            "unit": "ns/op",
            "extra": "364040192 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "364040192 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "364040192 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.711,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "452744240 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.711,
            "unit": "ns/op",
            "extra": "452744240 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "452744240 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "452744240 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.716,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "453764140 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.716,
            "unit": "ns/op",
            "extra": "453764140 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "453764140 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "453764140 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "442950828 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.77,
            "unit": "ns/op",
            "extra": "442950828 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "442950828 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "442950828 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.772,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "436767685 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.772,
            "unit": "ns/op",
            "extra": "436767685 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "436767685 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "436767685 times\n4 procs"
          }
        ]
      }
    ]
  }
}