package main

import (
	"atlas-fame/database"
	"atlas-fame/fame"
	fame2 "atlas-fame/kafka/consumer/fame"
	"atlas-fame/logger"
	"atlas-fame/service"
	"atlas-fame/tracing"
	"github.com/Chronicle20/atlas-kafka/consumer"
)

const serviceName = "atlas-fame"
const consumerGroupId = "Fame Service"

func main() {
	l := logger.CreateLogger(serviceName)
	l.Infoln("Starting main service.")

	tdm := service.GetTeardownManager()

	tc, err := tracing.InitTracer(l)(serviceName)
	if err != nil {
		l.WithError(err).Fatal("Unable to initialize tracer.")
	}

	db := database.Connect(l, database.SetMigrations(fame.Migration))

	cmf := consumer.GetManager().AddConsumer(l, tdm.Context(), tdm.WaitGroup())
	fame2.InitConsumers(l)(cmf)(consumerGroupId)
	fame2.InitHandlers(l)(db)(consumer.GetManager().RegisterHandler)

	tdm.TeardownFunc(tracing.Teardown(l)(tc))

	tdm.Wait()
	l.Infoln("Service shutdown.")
}
