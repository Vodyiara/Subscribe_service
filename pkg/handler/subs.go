package handler

import (
	"Test_project_Effective_Mobile"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// @Summary Создать подписку
// @Description Создает новую подписку с указанием сервиса, цены, start_date, end_date и user_id
// @Tags подписки
// @Accept json
// @Produce json
// @Param subscription body Test_project_Effective_Mobile.Subscripe true "Данные подписки"
// @Success 200 {object} map[string]int "Возвращает ID созданной подписки"
// @Failure 400 {object} errorResponse "Неверное тело запроса"
// @Failure 500 {object} errorResponse "Ошибка при создании подписки"
// @Router /subs [post]
func (h *Handler) createSubs(c *gin.Context) {

	var input Test_project_Effective_Mobile.Subscripe
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Warn("failed to bind JSON in createSubs")
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.service.TodoSub.Create(input)
	if err != nil {
		logrus.WithError(err).Errorf("failed to create subscription for user_id=%s service=%s", input.User_id, input.Service_name)
		newErrorResponse(c, http.StatusInternalServerError, "failed to create subscription")
		return
	}

	logrus.WithFields(logrus.Fields{
		"subscription_id": id,
		"user_id":         input.User_id,
		"service":         input.Service_name,
	}).Info("subscription created successfully")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllSubsResponse struct {
	Data []Test_project_Effective_Mobile.Subscripe "json:'data'"
}

// @Summary Получить все подписки
// @Description Возвращает список всех подписок
// @Tags подписки
// @Produce json
// @Success 200 {object} getAllSubsResponse "Список подписок"
// @Failure 500 {object} errorResponse "Ошибка при получении подписок"
// @Router /subs [get]
func (h *Handler) getAllSubs(c *gin.Context) {
	lists, err := h.service.TodoSub.GetAllSubs()
	if err != nil {
		logrus.WithError(err).Error("failed to get all subscriptions")
		newErrorResponse(c, http.StatusInternalServerError, "failed to fetch subscriptions")
		return
	}

	logrus.Infof("fetched %d subscriptions", len(lists))
	c.JSON(http.StatusOK, getAllSubsResponse{Data: lists})
}

// @Summary Обновить подписку
// @Description Обновляет поля подписки (service_name, price, start_date, end_date, done)
// @Tags подписки
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body Test_project_Effective_Mobile.UpdateSubInput true "Поля для обновления"
// @Success 200 {object} statusResponse "Статус Ok"
// @Failure 400 {object} errorResponse "Неверные данные или ID"
// @Failure 500 {object} errorResponse "Ошибка при обновлении подписки"
// @Router /subs/{id} [put]
func (h *Handler) updateSub(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithField("param", c.Param("id")).Warn("invalid id parameter in updateSub")
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}

	var input Test_project_Effective_Mobile.UpdateSubInput
	if err := c.BindJSON(&input); err != nil {
		logrus.WithError(err).Warn("failed to bind JSON in updateSub")
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.TodoSub.Update(id, input); err != nil {
		logrus.WithError(err).Errorf("failed to update subscription id=%d", id)
		newErrorResponse(c, http.StatusInternalServerError, "failed to update subscription")
		return
	}

	logrus.WithFields(logrus.Fields{
		"subscription_id": id,
		"service_name":    input.Service_name,
	}).Info("subscription updated successfully")

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags подписки
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} statusResponse "Статус Ok"
// @Failure 400 {object} errorResponse "Неверный ID"
// @Failure 500 {object} errorResponse "Ошибка при удалении подписки"
// @Router /subs/{id} [delete]
func (h *Handler) deleteSub(c *gin.Context) {
	subid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithField("param", c.Param("id")).Warn("invalid id parameter in deleteSub")
		newErrorResponse(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	if err := h.service.TodoSub.Delete(subid); err != nil {
		logrus.WithError(err).Errorf("failed to delete subscription id=%d", subid)
		newErrorResponse(c, http.StatusInternalServerError, "failed to delete subscription")
		return
	}

	logrus.WithField("subscription_id", subid).Info("subscription deleted successfully")
	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

// @Summary Получить подписку по ID
// @Description Возвращает детальную информацию о подписке по ID
// @Tags подписки
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} Test_project_Effective_Mobile.Subscripe "Данные подписки"
// @Failure 400 {object} errorResponse "Неверный ID"
// @Failure 500 {object} errorResponse "Ошибка при получении подписки"
// @Router /subs/{id} [get]
func (h *Handler) getSubById(c *gin.Context) {
	subid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.WithField("param", c.Param("id")).Warn("invalid id parameter in getSubById")
		newErrorResponse(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	sub, err := h.service.TodoSub.GetByIdSub(subid)
	if err != nil {
		logrus.WithError(err).Errorf("failed to get subscription id=%d", subid)
		newErrorResponse(c, http.StatusInternalServerError, "failed to fetch subscription")
		return
	}

	logrus.WithField("subscription_id", subid).Info("subscription fetched successfully")
	c.JSON(http.StatusOK, sub)
}

// @Summary Подсчитать общую стоимость подписок
// @Description Вычисляет сумму подписок с фильтром по user_id, service_name, start_date и end_date
// @Tags подписки
// @Produce json
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Param start_date query string true "Дата начала (MM-YYYY)"
// @Param end_date query string true "Дата окончания (MM-YYYY)"
// @Success 200 {object} map[string]int "Общая стоимость"
// @Failure 400 {object} errorResponse "Неверные параметры запроса"
// @Failure 500 {object} errorResponse "Ошибка при вычислении общей стоимости"
// @Router /subs/total [get]
func (h *Handler) getTotalPrice(c *gin.Context) {
	userIdStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		logrus.Warn("start_date or end_date not provided in getTotalPrice")
		newErrorResponse(c, http.StatusBadRequest, "start_date and end_date are required")
		return
	}

	var userId uuid.UUID
	var err error
	if userIdStr != "" {
		userId, err = uuid.Parse(userIdStr)
		if err != nil {
			logrus.WithField("user_id", userIdStr).Warn("invalid user_id format in getTotalPrice")
			newErrorResponse(c, http.StatusBadRequest, "invalid user_id format")
			return
		}
	}

	total, err := h.service.TodoSub.GetTotalPrice(userId, serviceName, startDate, endDate)
	if err != nil {
		logrus.WithError(err).Errorf("failed to calculate total price for user_id=%s service=%s", userId, serviceName)
		newErrorResponse(c, http.StatusInternalServerError, "failed to calculate total price")
		return
	}

	logrus.WithFields(logrus.Fields{
		"user_id":      userId,
		"service_name": serviceName,
		"start_date":   startDate,
		"end_date":     endDate,
		"total_price":  total,
	}).Info("calculated total subscription price")

	c.JSON(http.StatusOK, gin.H{"total_price": total})
}
