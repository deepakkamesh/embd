package embd

import (
	"fmt"

	"github.com/golang/glog"
)

type pin interface {
	Close() error
}

type gpioDriver struct {
	pinMap          PinMap
	initializedPins map[string]pin
}

func newGPIODriver(pinMap PinMap) *gpioDriver {
	return &gpioDriver{
		pinMap:          pinMap,
		initializedPins: map[string]pin{},
	}
}

func (io *gpioDriver) lookupKey(key interface{}, cap int) (*PinDesc, bool) {
	return io.pinMap.Lookup(key, cap)
}

func (io *gpioDriver) digitalPin(key interface{}) (*digitalPin, error) {
	pd, found := io.lookupKey(key, CapNormal)
	if !found {
		return nil, fmt.Errorf("gpio: could not find pin matching %q", key)
	}

	if pd.Caps != CapNormal {
		glog.Infof("gpio: pin %q is not a dedicated digital io pin. please refer to the system reference manual for more details", key)
	}

	dp := newDigitalPin(pd.DigitalLogical)
	io.initializedPins[pd.ID] = dp

	return dp, nil
}

func (io *gpioDriver) DigitalPin(key interface{}) (DigitalPin, error) {
	return io.digitalPin(key)
}

func (io *gpioDriver) Close() error {
	for _, p := range io.initializedPins {
		if err := p.Close(); err != nil {
			return err
		}
	}

	return nil
}
