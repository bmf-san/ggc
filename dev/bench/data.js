window.BENCHMARK_DATA = {
  "lastUpdate": 1778852564333,
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
          "id": "d5143a4f0cc33077328e325b699f7f4270801f02",
          "message": "chore(deps): bump golang.org/x/text from 0.36.0 to 0.37.0 (#424)\n\nBumps [golang.org/x/text](https://github.com/golang/text) from 0.36.0 to 0.37.0.\n- [Release notes](https://github.com/golang/text/releases)\n- [Commits](https://github.com/golang/text/compare/v0.36.0...v0.37.0)\n\n---\nupdated-dependencies:\n- dependency-name: golang.org/x/text\n  dependency-version: 0.37.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-05-14T22:50:52+09:00",
          "tree_id": "92c81435c6f052e273e8ff2fa9a119ae765ed2c5",
          "url": "https://github.com/bmf-san/ggc/commit/d5143a4f0cc33077328e325b699f7f4270801f02"
        },
        "date": 1778766844298,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22179624,
            "unit": "ns/op\t 1382776 B/op\t    6920 allocs/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22179624,
            "unit": "ns/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1382776,
            "unit": "B/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6920,
            "unit": "allocs/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22196053,
            "unit": "ns/op\t 1385314 B/op\t    6920 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22196053,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1385314,
            "unit": "B/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6920,
            "unit": "allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22294635,
            "unit": "ns/op\t 1390505 B/op\t    6920 allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22294635,
            "unit": "ns/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1390505,
            "unit": "B/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6920,
            "unit": "allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22380514,
            "unit": "ns/op\t 1394637 B/op\t    6922 allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22380514,
            "unit": "ns/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1394637,
            "unit": "B/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6922,
            "unit": "allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22172176,
            "unit": "ns/op\t 1390746 B/op\t    6920 allocs/op",
            "extra": "60 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22172176,
            "unit": "ns/op",
            "extra": "60 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1390746,
            "unit": "B/op",
            "extra": "60 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6920,
            "unit": "allocs/op",
            "extra": "60 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22212941,
            "unit": "ns/op\t 1396825 B/op\t    6921 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22212941,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1396825,
            "unit": "B/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 6921,
            "unit": "allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22395704,
            "unit": "ns/op\t 1475111 B/op\t    7978 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22395704,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1475111,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7978,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22359693,
            "unit": "ns/op\t 1474379 B/op\t    7977 allocs/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22359693,
            "unit": "ns/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474379,
            "unit": "B/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7977,
            "unit": "allocs/op",
            "extra": "97 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22310750,
            "unit": "ns/op\t 1474883 B/op\t    7978 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22310750,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1474883,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7978,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22404845,
            "unit": "ns/op\t 1475200 B/op\t    7978 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22404845,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1475200,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7978,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22406854,
            "unit": "ns/op\t 1475797 B/op\t    7979 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22406854,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1475797,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7979,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22399673,
            "unit": "ns/op\t 1475077 B/op\t    7978 allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22399673,
            "unit": "ns/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1475077,
            "unit": "B/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7978,
            "unit": "allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3083,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "382401 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3083,
            "unit": "ns/op",
            "extra": "382401 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "382401 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "382401 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3172,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "374126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3172,
            "unit": "ns/op",
            "extra": "374126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "374126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "374126 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3086,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "374796 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3086,
            "unit": "ns/op",
            "extra": "374796 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "374796 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "374796 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3081,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "375535 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3081,
            "unit": "ns/op",
            "extra": "375535 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "375535 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "375535 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3097,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "369124 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3097,
            "unit": "ns/op",
            "extra": "369124 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "369124 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "369124 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3127,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "368815 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3127,
            "unit": "ns/op",
            "extra": "368815 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "368815 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "368815 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 66.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17998208 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 66.31,
            "unit": "ns/op",
            "extra": "17998208 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17998208 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17998208 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18285602 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.71,
            "unit": "ns/op",
            "extra": "18285602 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18285602 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18285602 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 65.72,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18227437 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 65.72,
            "unit": "ns/op",
            "extra": "18227437 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18227437 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18227437 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18597342 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 67,
            "unit": "ns/op",
            "extra": "18597342 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18597342 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18597342 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 65.16,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18264490 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 65.16,
            "unit": "ns/op",
            "extra": "18264490 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18264490 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18264490 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18496596 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.5,
            "unit": "ns/op",
            "extra": "18496596 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18496596 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18496596 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6617,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "179898 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6617,
            "unit": "ns/op",
            "extra": "179898 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "179898 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "179898 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6620,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "177681 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6620,
            "unit": "ns/op",
            "extra": "177681 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "177681 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "177681 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6602,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "177064 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6602,
            "unit": "ns/op",
            "extra": "177064 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "177064 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "177064 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6605,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "178274 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6605,
            "unit": "ns/op",
            "extra": "178274 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "178274 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "178274 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6610,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "177694 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6610,
            "unit": "ns/op",
            "extra": "177694 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "177694 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "177694 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6817,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "178140 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6817,
            "unit": "ns/op",
            "extra": "178140 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "178140 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "178140 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1570,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "753526 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1570,
            "unit": "ns/op",
            "extra": "753526 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "753526 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "753526 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1562,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "753567 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1562,
            "unit": "ns/op",
            "extra": "753567 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "753567 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "753567 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1560,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "752163 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1560,
            "unit": "ns/op",
            "extra": "752163 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "752163 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "752163 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1567,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "755950 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1567,
            "unit": "ns/op",
            "extra": "755950 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "755950 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "755950 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1559,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "735058 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1559,
            "unit": "ns/op",
            "extra": "735058 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "735058 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "735058 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1573,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "762193 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1573,
            "unit": "ns/op",
            "extra": "762193 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "762193 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "762193 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1114,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1114,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1107,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1107,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1115,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1115,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1104,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1104,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1106,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1106,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1108,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1108,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1203,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "999624 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1203,
            "unit": "ns/op",
            "extra": "999624 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "999624 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "999624 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1199,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "997086 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1199,
            "unit": "ns/op",
            "extra": "997086 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "997086 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "997086 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1201,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1201,
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
            "value": 1199,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1199,
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
            "value": 1199,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "995798 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1199,
            "unit": "ns/op",
            "extra": "995798 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "995798 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "995798 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1200,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1200,
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
            "value": 1589,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "782252 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1589,
            "unit": "ns/op",
            "extra": "782252 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "782252 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "782252 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1584,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "778389 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1584,
            "unit": "ns/op",
            "extra": "778389 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "778389 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "778389 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1585,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "772184 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1585,
            "unit": "ns/op",
            "extra": "772184 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "772184 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "772184 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1584,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "759376 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1584,
            "unit": "ns/op",
            "extra": "759376 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "759376 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "759376 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1583,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "766400 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1583,
            "unit": "ns/op",
            "extra": "766400 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "766400 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "766400 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1585,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "759710 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1585,
            "unit": "ns/op",
            "extra": "759710 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "759710 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "759710 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1100,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1100,
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
            "value": 1125,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1125,
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
            "value": 1106,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1106,
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
            "value": 1126,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1126,
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
            "value": 1126,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1126,
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
            "value": 1086,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1086,
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
            "value": 1077,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1077,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1077,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1077,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1079,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1079,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1081,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1081,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1080,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1080,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1078,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1078,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1306,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "964116 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1306,
            "unit": "ns/op",
            "extra": "964116 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "964116 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "964116 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1249,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "919716 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1249,
            "unit": "ns/op",
            "extra": "919716 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "919716 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "919716 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1253,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "937749 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1253,
            "unit": "ns/op",
            "extra": "937749 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "937749 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "937749 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1255,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "957592 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1255,
            "unit": "ns/op",
            "extra": "957592 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "957592 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "957592 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1249,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "958005 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1249,
            "unit": "ns/op",
            "extra": "958005 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "958005 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "958005 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1253,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "965815 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1253,
            "unit": "ns/op",
            "extra": "965815 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "965815 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "965815 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 3.093,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "330307971 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 3.093,
            "unit": "ns/op",
            "extra": "330307971 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "330307971 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "330307971 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.916,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "427892182 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.916,
            "unit": "ns/op",
            "extra": "427892182 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "427892182 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "427892182 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.909,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "426660927 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.909,
            "unit": "ns/op",
            "extra": "426660927 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "426660927 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "426660927 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.908,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "416728258 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.908,
            "unit": "ns/op",
            "extra": "416728258 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "416728258 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "416728258 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.876,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "413562847 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.876,
            "unit": "ns/op",
            "extra": "413562847 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "413562847 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "413562847 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.892,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "425513997 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.892,
            "unit": "ns/op",
            "extra": "425513997 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "425513997 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "425513997 times\n4 procs"
          }
        ]
      },
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
          "id": "459882b869adf2bcdd86adbaae1d845763daac42",
          "message": "feat(cmd): add 'show' command (#430)\n\nImplements 'ggc show' as a wrapper around 'git show' for inspecting\ncommits, tags, trees, and blobs. Supports passthrough of standard\ngit show options (e.g. --stat, --name-only) and accepts any number\nof object arguments.\n\nrefs #428",
          "timestamp": "2026-05-15T22:39:27+09:00",
          "tree_id": "7d6c2034685dffa7a8fb5726ca79b4d1a7519ff5",
          "url": "https://github.com/bmf-san/ggc/commit/459882b869adf2bcdd86adbaae1d845763daac42"
        },
        "date": 1778852563674,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22909252,
            "unit": "ns/op\t 1433583 B/op\t    7296 allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22909252,
            "unit": "ns/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1433583,
            "unit": "B/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7296,
            "unit": "allocs/op",
            "extra": "48 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22812644,
            "unit": "ns/op\t 1432017 B/op\t    7295 allocs/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22812644,
            "unit": "ns/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1432017,
            "unit": "B/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7295,
            "unit": "allocs/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22514451,
            "unit": "ns/op\t 1423862 B/op\t    7294 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22514451,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1423862,
            "unit": "B/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7294,
            "unit": "allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22571019,
            "unit": "ns/op\t 1432206 B/op\t    7295 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22571019,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1432206,
            "unit": "B/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7295,
            "unit": "allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22651192,
            "unit": "ns/op\t 1434600 B/op\t    7296 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22651192,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1434600,
            "unit": "B/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7296,
            "unit": "allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8)",
            "value": 22671160,
            "unit": "ns/op\t 1433323 B/op\t    7294 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22671160,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1433323,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Version (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 7294,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22728086,
            "unit": "ns/op\t 1536375 B/op\t    8390 allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22728086,
            "unit": "ns/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1536375,
            "unit": "B/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 8390,
            "unit": "allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22767554,
            "unit": "ns/op\t 1536444 B/op\t    8390 allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22767554,
            "unit": "ns/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1536444,
            "unit": "B/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 8390,
            "unit": "allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 23009115,
            "unit": "ns/op\t 1535249 B/op\t    8390 allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 23009115,
            "unit": "ns/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1535249,
            "unit": "B/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 8390,
            "unit": "allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22737259,
            "unit": "ns/op\t 1536567 B/op\t    8390 allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22737259,
            "unit": "ns/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1536567,
            "unit": "B/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 8390,
            "unit": "allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 30740717,
            "unit": "ns/op\t 1536528 B/op\t    8391 allocs/op",
            "extra": "96 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 30740717,
            "unit": "ns/op",
            "extra": "96 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1536528,
            "unit": "B/op",
            "extra": "96 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 8391,
            "unit": "allocs/op",
            "extra": "96 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8)",
            "value": 22816684,
            "unit": "ns/op\t 1536984 B/op\t    8391 allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - ns/op",
            "value": 22816684,
            "unit": "ns/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - B/op",
            "value": 1536984,
            "unit": "B/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkStartup_Help (github.com/bmf-san/ggc/v8) - allocs/op",
            "value": 8391,
            "unit": "allocs/op",
            "extra": "98 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3110,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "379014 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3110,
            "unit": "ns/op",
            "extra": "379014 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "379014 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "379014 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3110,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "371250 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3110,
            "unit": "ns/op",
            "extra": "371250 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "371250 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "371250 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3094,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "369205 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3094,
            "unit": "ns/op",
            "extra": "369205 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "369205 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "369205 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3092,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "367233 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3092,
            "unit": "ns/op",
            "extra": "367233 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "367233 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "367233 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3108,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "372408 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3108,
            "unit": "ns/op",
            "extra": "372408 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "372408 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "372408 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd)",
            "value": 3103,
            "unit": "ns/op\t    1032 B/op\t      29 allocs/op",
            "extra": "368518 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 3103,
            "unit": "ns/op",
            "extra": "368518 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 1032,
            "unit": "B/op",
            "extra": "368518 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Default (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 29,
            "unit": "allocs/op",
            "extra": "368518 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18504116 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.38,
            "unit": "ns/op",
            "extra": "18504116 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18504116 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18504116 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18673333 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.2,
            "unit": "ns/op",
            "extra": "18673333 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18673333 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18673333 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18407430 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.44,
            "unit": "ns/op",
            "extra": "18407430 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18407430 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18407430 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18689493 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.78,
            "unit": "ns/op",
            "extra": "18689493 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18689493 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18689493 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18443949 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.99,
            "unit": "ns/op",
            "extra": "18443949 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18443949 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18443949 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd)",
            "value": 64.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18479463 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - ns/op",
            "value": 64.97,
            "unit": "ns/op",
            "extra": "18479463 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18479463 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugger_DebugKeys_Help (github.com/bmf-san/ggc/v8/cmd) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18479463 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6630,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "181014 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6630,
            "unit": "ns/op",
            "extra": "181014 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "181014 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "181014 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6621,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "176628 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6621,
            "unit": "ns/op",
            "extra": "176628 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "176628 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "176628 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6692,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "175314 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6692,
            "unit": "ns/op",
            "extra": "175314 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "175314 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "175314 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6626,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "176187 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6626,
            "unit": "ns/op",
            "extra": "176187 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "176187 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "176187 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6635,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "177378 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6635,
            "unit": "ns/op",
            "extra": "177378 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "177378 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "177378 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 6632,
            "unit": "ns/op\t     961 B/op\t      48 allocs/op",
            "extra": "174804 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 6632,
            "unit": "ns/op",
            "extra": "174804 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 961,
            "unit": "B/op",
            "extra": "174804 times\n4 procs"
          },
          {
            "name": "BenchmarkManagerGet (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 48,
            "unit": "allocs/op",
            "extra": "174804 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1593,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "753984 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1593,
            "unit": "ns/op",
            "extra": "753984 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "753984 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "753984 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1570,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "758284 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1570,
            "unit": "ns/op",
            "extra": "758284 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "758284 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "758284 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1575,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "765396 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1575,
            "unit": "ns/op",
            "extra": "765396 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "765396 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "765396 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1573,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "757748 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1573,
            "unit": "ns/op",
            "extra": "757748 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "757748 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "757748 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1568,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "747274 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1568,
            "unit": "ns/op",
            "extra": "747274 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "747274 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "747274 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config)",
            "value": 1571,
            "unit": "ns/op\t     128 B/op\t       3 allocs/op",
            "extra": "757808 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - ns/op",
            "value": 1571,
            "unit": "ns/op",
            "extra": "757808 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - B/op",
            "value": 128,
            "unit": "B/op",
            "extra": "757808 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizeConfigPath (github.com/bmf-san/ggc/v8/internal/config) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "757808 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1105,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1105,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1105,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1105,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1108,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1108,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1108,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1108,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1105,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1105,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1107,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1107,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/short (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1199,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "999218 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1199,
            "unit": "ns/op",
            "extra": "999218 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "999218 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "999218 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1204,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "996387 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1204,
            "unit": "ns/op",
            "extra": "996387 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "996387 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "996387 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1201,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1201,
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
            "value": 1203,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1203,
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
            "value": 1199,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "964858 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1199,
            "unit": "ns/op",
            "extra": "964858 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "964858 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "964858 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1201,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/medium (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1201,
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
            "value": 1580,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "748243 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1580,
            "unit": "ns/op",
            "extra": "748243 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "748243 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "748243 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1700,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "774169 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1700,
            "unit": "ns/op",
            "extra": "774169 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "774169 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "774169 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1581,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "770923 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1581,
            "unit": "ns/op",
            "extra": "770923 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "770923 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "770923 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1583,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "761534 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1583,
            "unit": "ns/op",
            "extra": "761534 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "761534 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "761534 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1585,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "729447 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1585,
            "unit": "ns/op",
            "extra": "729447 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "729447 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "729447 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1596,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "768720 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1596,
            "unit": "ns/op",
            "extra": "768720 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "768720 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/long (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "768720 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1117,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1117,
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
            "value": 1149,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1149,
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
            "value": 1154,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "911570 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1154,
            "unit": "ns/op",
            "extra": "911570 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "911570 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "911570 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1154,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1154,
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
            "value": 1131,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1131,
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
            "value": 1116,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/typo_miss (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1116,
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
            "value": 1075,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1107906 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1075,
            "unit": "ns/op",
            "extra": "1107906 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1107906 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1107906 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1104,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1104,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1091,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1091,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1088,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1088,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1084,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1084,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1084,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1084,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatch/single_char (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1253,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "949250 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1253,
            "unit": "ns/op",
            "extra": "949250 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "949250 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "949250 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1259,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "944842 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1259,
            "unit": "ns/op",
            "extra": "944842 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "944842 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "944842 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1258,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "960543 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1258,
            "unit": "ns/op",
            "extra": "960543 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "960543 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "960543 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1250,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "959493 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1250,
            "unit": "ns/op",
            "extra": "959493 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "959493 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "959493 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1261,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "945560 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1261,
            "unit": "ns/op",
            "extra": "945560 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "945560 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "945560 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive)",
            "value": 1253,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "934012 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - ns/op",
            "value": 1253,
            "unit": "ns/op",
            "extra": "934012 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "934012 times\n4 procs"
          },
          {
            "name": "BenchmarkFuzzyMatchScore (github.com/bmf-san/ggc/v8/internal/interactive) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "934012 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.937,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "421205745 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.937,
            "unit": "ns/op",
            "extra": "421205745 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "421205745 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "421205745 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.964,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "408507912 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.964,
            "unit": "ns/op",
            "extra": "408507912 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "408507912 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "408507912 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "416385422 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.83,
            "unit": "ns/op",
            "extra": "416385422 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "416385422 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "416385422 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.899,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "420605019 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.899,
            "unit": "ns/op",
            "extra": "420605019 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "420605019 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "420605019 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 2.857,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "417607365 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 2.857,
            "unit": "ns/op",
            "extra": "417607365 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "417607365 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "417607365 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings)",
            "value": 3.002,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "413749154 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - ns/op",
            "value": 3.002,
            "unit": "ns/op",
            "extra": "413749154 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "413749154 times\n4 procs"
          },
          {
            "name": "BenchmarkDebugKeysCommand_IdentifySequence (github.com/bmf-san/ggc/v8/internal/keybindings) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "413749154 times\n4 procs"
          }
        ]
      }
    ]
  }
}