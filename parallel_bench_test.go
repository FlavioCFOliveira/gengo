package gengo

// Parallel benchmarks for task #39 (contention under concurrent use).
//
// Each BenchmarkParallel<Func> drives a package-level generator from
// GOMAXPROCS goroutines through b.RunParallel. Run them with -cpu to measure
// scaling, for example:
//
//	go test -run '^$' -bench '^BenchmarkParallel' -benchmem -cpu 1,2,4,8,16 .
//
// Every goroutine folds each result into a goroutine-local accumulator and
// publishes it once, atomically, when its loop ends. This keeps the compiler
// from discarding the call or stack-allocating a result the caller would
// normally retain, and it adds no shared write to the measured loop: a shared
// per-iteration sink would itself be a contention source (cache-line
// ping-pong) and would bias the measurement.

import (
	"math"
	"sync/atomic"
	"testing"
	"time"
)

// benchParallelSink receives each goroutine's accumulator once per benchmark.
var benchParallelSink atomic.Uint64

// benchParallelStart and benchParallelEnd bound the DateBetween benchmark.
var (
	benchParallelStart = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	benchParallelEnd   = time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC)
)

// ---------- Numbers ----------

func BenchmarkParallelInt8(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int8())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt8Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int8Between(10, 127))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt16(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int16())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt16Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int16Between(-1000, 1000))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt32(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int32())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt32Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int32Between(-100000, 100000))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelIntBetween(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(IntBetween(-100000, 100000))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt64(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int64())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInt64Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Int64Between(-1000000000, 1000000000))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelByte(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Byte())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint8(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Uint8())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint8Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Uint8Between(10, 200))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint16(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Uint16())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint16Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Uint16Between(10, 60000))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint32(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Uint32())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint32Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Uint32Between(10, 4000000000))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint64(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += Uint64()
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUint64Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += Uint64Between(10, 1<<62)
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelFloat32(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(math.Float32bits(Float32()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelFloat32Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(math.Float32bits(Float32Between(-1000, 1000)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelFloat64(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += math.Float64bits(Float64())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelFloat64Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += math.Float64bits(Float64Between(-1000, 1000))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelComplex64(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			v := Complex64()
			acc += uint64(math.Float32bits(real(v)) ^ math.Float32bits(imag(v)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelComplex64Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			v := Complex64Between(-10, 10, -10, 10)
			acc += uint64(math.Float32bits(real(v)) ^ math.Float32bits(imag(v)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelComplex128(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			v := Complex128()
			acc += math.Float64bits(real(v)) ^ math.Float64bits(imag(v))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelComplex128Between(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			v := Complex128Between(-10, 10, -10, 10)
			acc += math.Float64bits(real(v)) ^ math.Float64bits(imag(v))
		}
		benchParallelSink.Add(acc)
	})
}

// ---------- Bool ----------

func BenchmarkParallelBool(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			if Bool() {
				acc++
			}
		}
		benchParallelSink.Add(acc)
	})
}

// ---------- Dates ----------

func BenchmarkParallelDate(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(Date().Unix())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelUnixDate(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(UnixDate().Unix())
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelDateBetween(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(DateBetween(benchParallelStart, benchParallelEnd).Unix())
		}
		benchParallelSink.Add(acc)
	})
}

// ---------- Strings ----------

func BenchmarkParallelString(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(String(8, Alphanumeric)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelString64(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(String(64, Alphanumeric)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringBetween(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringBetween(4, 16, Alphanumeric)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringAllChars(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringAllChars(8)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringAlphanumeric(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringAlphanumeric(8)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringAlphabetic(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringAlphabetic(8)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringAlphabeticUppercase(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringAlphabeticUppercase(8)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringAlphabeticLowercase(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringAlphabeticLowercase(8)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringNumeric(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringNumeric(8)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringHexadecimal(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringHexadecimal(8)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelStringSymbols(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(StringSymbols(8)))
		}
		benchParallelSink.Add(acc)
	})
}

// ---------- Words ----------

func BenchmarkParallelWord(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(Word()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelWordByLengthTypeAny(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(WordByLengthType(AnyLengthWord)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelWordByLengthTypeSmall(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(WordByLengthType(SmallLengthWord)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelWordByLengthTypeMedium(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(WordByLengthType(MediumLengthWords)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelWordByLengthTypeBig(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(WordByLengthType(BigLengthWords)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelWords(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(Words(10)))
		}
		benchParallelSink.Add(acc)
	})
}

// ---------- WordsPT ----------

func BenchmarkParallelWordPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(WordPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelWordPTByLengthType(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(WordPTByLengthType(MediumLengthWords)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelNounPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(NounPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelNounPTOf(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(NounPTOf(Feminine, Plural, MediumLengthWords)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelAdjectivePT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(AdjectivePT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelAdjectivePTOf(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(AdjectivePTOf(Masculine, Plural, Superlative, AnyLengthWord)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelVerbPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(VerbPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelVerbPTOf(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(VerbPTOf(Indicative, Present, Third, Plural, AnyLengthWord)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelAdverbPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(AdverbPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelAdverbPTByLengthType(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(AdverbPTByLengthType(BigLengthWords)))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelArticlePT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(ArticlePT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelPrepositionPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(PrepositionPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelConjunctionPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(ConjunctionPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelPronounPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(PronounPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelInterjectionPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(InterjectionPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelNumeralPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(NumeralPT()))
		}
		benchParallelSink.Add(acc)
	})
}

func BenchmarkParallelWordsPT(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var acc uint64
		for pb.Next() {
			acc += uint64(len(WordsPT(16)))
		}
		benchParallelSink.Add(acc)
	})
}
