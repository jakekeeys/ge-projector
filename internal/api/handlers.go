package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jakekeeys/givforecast/internal/solcast"
)

func (s *Server) UpdateChargeTargetHandler(c *gin.Context) {
	err := s.UpdateChargeTarget()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *Server) ForecastNowHandler(c *gin.Context) {
	t := time.Now().Local()

	ts := c.Query("time")
	if ts != "" {
		pt, err := time.ParseInLocation(timeFormat, ts, time.Local)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		t = pt
	}

	fn, err := s.f.ForecastNow(t)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fn)
}

func (s *Server) ForecastHandler(c *gin.Context) {
	ds := c.Query("date")
	if ds == "" {
		ds = time.Now().Local().Format(dateFormat)
	}

	d, err := time.Parse(dateFormat, ds)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	today := time.Now().Local().Truncate(time.Hour * 24)
	if d.Before(today) {
		c.String(http.StatusBadRequest, "date must be today or < 7 days in the future")
		return
	}

	if d.After(today.AddDate(0, 0, 6)) {
		c.String(http.StatusBadRequest, "date must be today or < 7 days in the future")
		return
	}

	fc, err := s.f.Forecast(d)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, fc)
}

func (s *Server) ConfigHandler(c *gin.Context) {
	config := s.f.GetConfig()

	c.JSON(http.StatusOK, config)
}

func (s *Server) SetConsumptionAveragesHandler(c *gin.Context) {
	var data map[time.Time]float64
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	s.gec.SetConsumptionAverages(data)
	return
}

func (s *Server) GetConsumptionAveragesHandler(c *gin.Context) {
	averages, err := s.gec.GetConsumptionAverages()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, averages)
}

func (s *Server) UpdateConsumptionAveragesHandler(c *gin.Context) {
	err := s.gec.UpdateConsumptionAverages()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *Server) SetForecastDataHandler(c *gin.Context) {
	var data solcast.ForecastData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = s.sc.SetForecast(data)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	return
}

func (s *Server) UpdateForecastDataHandler(c *gin.Context) {
	err := s.sc.UpdateForecast()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	return
}