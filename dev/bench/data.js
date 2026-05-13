window.BENCHMARK_DATA = {
  "lastUpdate": 1778682229085,
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
      }
    ]
  }
}