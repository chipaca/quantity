package quantity

import (
	"fmt"
	"math"
)

// An Amount is an adimensional (i.e. unitless) countable quantity
type Amount uint64

func (n Amount) Format(f fmt.State, _ rune) {
	w, ok := f.Width()
	if !ok {
		w = -1
	}

	fmt.Fprint(f, FormatAmount(uint64(n), w))
}

// FormatAmount serializes an adimensional (i.e. unitless) positive
// integer into a given fixed width, using SI suffixes.
//
// With cannot be less than 3. If less than 0, a default width of 5 is
// used; if less than 3, 3 is used.
//
// The output string will be ascii right now, but this might change
// (there are some cases where using vulgar fractions would use the
// space a little better, e.g. FormatAmount(1500, 3) -> "1½k" instead
// of the current " 2k").
func FormatAmount(amount uint64, width int) string {
	if width < 0 {
		width = 5
	}
	max := uint64(5000)
	maxFloat := 999.5

	if width < 4 {
		width = 3
		max = 999
		maxFloat = 99.5
	}

	if amount <= max {
		pad := ""
		if width > 5 {
			pad = " "
		}
		return fmt.Sprintf("%*d%s", width-len(pad), amount, pad)
	}
	var prefix rune
	r := float64(amount)
	// zetta and yotta are me being pedantic: maxuint64 is ~18EB
	for _, prefix = range "kMGTPEZY" {
		r /= 1000
		if r < maxFloat {
			break
		}
	}

	width--
	digits := 3
	if r < 99.5 {
		digits--
		if r < 9.5 {
			digits--
			if r < .95 {
				digits--
			}
		}
	}
	precision := 0
	if (width - digits) > 1 {
		precision = width - digits - 1
	}

	s := fmt.Sprintf("%*.*f%c", width, precision, r, prefix)
	if r < .95 {
		return s[1:]
	}
	return s
}

// Bytes is an amount of bytes
type Bytes uint64

func (n Bytes) Format(f fmt.State, _ rune) {
	w, ok := f.Width()
	if !ok {
		w = -1
	}

	fmt.Fprint(f, FormatBytes(uint64(n), w))
}

// FormatBytes calls FormatAmount, and adds a 'B'. High tech stuff.
func FormatBytes(n uint64, width int) string {
	return FormatAmount(n, width-1) + "B"
}

// BPS are bytes per second, a bitrate.
type BPS struct {
	Bytes    Bytes
	Duration Duration
}

func (n BPS) Format(f fmt.State, _ rune) {
	w, ok := f.Width()
	if !ok {
		w = -1
	}

	fmt.Fprint(f, FormatBPS(n.Bytes, n.Duration, w))
}

// FormatBPS formats a bitrate. Minimum width is 6, defaults is 8. dt
// must be positive.
func FormatBPS(n Bytes, sec Duration, width int) string {
	if sec < 0 {
		sec = -sec
	}
	return FormatBytes(uint64(float64(n)/float64(sec)), width-2) + "/s"
}

// Duration in fractions of a second
type Duration float64

func (dt Duration) Format(f fmt.State, _ rune) {
	fmt.Fprint(f, FormatDuration(dt))
}

const (
	period = 365.25 // julian years (c.f. the actual orbital period, 365.256363004d)
)

func divmod(a, b float64) (q, r float64) {
	q = math.Floor(a / b)
	return q, a - q*b
}

// FormatDuration formats a Duration into a string of length 5.
//
// TODO: make the width a parameter. Also look into being smarter;
// this is a rather crude and brute-forceish (a.k.a. forest of ifs)
// approach. OTOH it really is a lot of special cases...
func FormatDuration(dur Duration) string {
	dt := float64(dur)
	if dt < 60 {
		if dt >= 9.995 {
			return fmt.Sprintf("%.1fs", dt)
		} else if dt >= .9995 {
			return fmt.Sprintf("%.2fs", dt)
		}

		var prefix rune
		for _, prefix = range "mun" {
			dt *= 1000
			if dt >= .9995 {
				break
			}
		}

		if dt > 9.5 {
			return fmt.Sprintf("%3.f%cs", dt, prefix)
		}

		return fmt.Sprintf("%.1f%cs", dt, prefix)
	}

	if dt < 600 {
		m, s := divmod(dt, 60)
		return fmt.Sprintf("%.fm%02.fs", m, s)
	}

	dt /= 60 // dt now minutes

	if dt < 99.95 {
		return fmt.Sprintf("%3.1fm", dt)
	}

	if dt < 10*60 {
		h, m := divmod(dt, 60)
		return fmt.Sprintf("%.fh%02.fm", h, m)
	}

	if dt < 24*60 {
		if h, m := divmod(dt, 60); m < 10 {
			return fmt.Sprintf("%.fh%1.fm", h, m)
		}

		return fmt.Sprintf("%3.1fh", dt/60)
	}

	dt /= 60 // dt now hours

	if dt < 10*24 {
		d, h := divmod(dt, 24)
		return fmt.Sprintf("%.fd%02.fh", d, h)
	}

	if dt < 99.95*24 {
		if d, h := divmod(dt, 24); h < 10 {
			return fmt.Sprintf("%.fd%.fh", d, h)
		}
		return fmt.Sprintf("%4.1fd", dt/24)
	}

	dt /= 24 // dt now days

	if dt < 2*period {
		return fmt.Sprintf("%4.0fd", dt)
	}

	dt /= period // dt now years

	if dt < 9.995 {
		return fmt.Sprintf("%4.2fy", dt)
	}

	if dt < 99.95 {
		return fmt.Sprintf("%4.1fy", dt)
	}

	return fmt.Sprintf("%4.fy", dt)
}
