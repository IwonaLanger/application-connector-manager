package reconciler

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type stateFn func(context.Context, *fsm, *systemState) (stateFn, *ctrl.Result, error)

type Watch = func(src source.Source, eventhandler handler.EventHandler, predicates ...predicate.Predicate) error

type K8s struct {
	client.Client
	record.EventRecorder
	Watch
	handler.MapFunc
}

type Fsm interface {
	Run(ctx context.Context, v v1alpha1.ApplicationConnector) (ctrl.Result, error)
}

type fsm struct {
	fn  stateFn
	log *zap.SugaredLogger
	K8s
	Cfg
	dependencyACK *bool
}

func (m *fsm) stateFnName() string {
	fullName := runtime.FuncForPC(reflect.ValueOf(m.fn).Pointer()).Name()
	splitFullName := strings.Split(fullName, ".")

	if len(splitFullName) < 3 {
		return fullName
	}

	shortName := splitFullName[2]
	return shortName
}

var mux sync.Mutex

func (m *fsm) Run(ctx context.Context, v v1alpha1.ApplicationConnector) (ctrl.Result, error) {
	state := systemState{Instance: v}
	var err error
	var result *ctrl.Result
	mux.Lock()
loop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break loop
		default:
			stateFnName := m.stateFnName()
			m.fn, result, err = m.fn(ctx, m, &state)
			newStateFnName := m.stateFnName()
			m.log.With("result", result, "err", err, "mFnIsNill", m.fn == nil).Info(fmt.Sprintf("switching state from %s to %s", stateFnName, newStateFnName))
			if m.fn == nil || err != nil {
				break loop
			}
		}
	}
	mux.Unlock()

	m.log.With("error", err).
		With("result", result).
		Info("reconciliation done")

	if result != nil {
		return *result, err
	}

	return ctrl.Result{
		Requeue: false,
	}, err
}

func NewFsm(log *zap.SugaredLogger, cfg Cfg, k8s K8s, depsACK *bool) Fsm {
	return &fsm{
		fn:            sFnTakeSnapshot,
		Cfg:           cfg,
		log:           log,
		K8s:           k8s,
		dependencyACK: depsACK,
	}
}
