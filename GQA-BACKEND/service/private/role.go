package private

import (
	"errors"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/global"
	"github.com/Junvary/gin-quasar-admin/GQA-BACKEND/model"
	"gorm.io/gorm"
)

type ServiceRole struct{}

func (s *ServiceRole) GetRoleList(requestRoleList model.RequestGetRoleList) (err error, role interface{}, total int64) {
	pageSize := requestRoleList.PageSize
	offset := requestRoleList.PageSize * (requestRoleList.Page - 1)
	db := global.GqaDb.Model(&model.SysRole{})
	var roleList []model.SysRole
	//配置搜索
	if requestRoleList.RoleCode != "" {
		db = db.Where("role_code like ?", "%"+requestRoleList.RoleCode+"%")
	}
	if requestRoleList.RoleName != "" {
		db = db.Where("role_name like ?", "%"+requestRoleList.RoleName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(pageSize).Offset(offset).Order(model.OrderByColumn(requestRoleList.SortBy, requestRoleList.Desc)).
		Preload("DefaultPageMenu").Find(&roleList).Error
	return err, roleList, total
}

func (s *ServiceRole) EditRole(toEditRole model.SysRole) (err error) {
	var sysRole model.SysRole
	if sysRole.Stable == "yesNo_yes" {
		return errors.New("系统内置不允许编辑：" + toEditRole.RoleCode)
	}
	if err = global.GqaDb.Where("id = ?", toEditRole.Id).First(&sysRole).Error; err != nil {
		return err
	}
	//err = global.GqaDb.Updates(&toEditRole).Error
	err = global.GqaDb.Save(&toEditRole).Error
	return err
}

func (s *ServiceRole) AddRole(toAddRole model.SysRole) (err error) {
	var role model.SysRole
	if !errors.Is(global.GqaDb.Where("role_code = ?", toAddRole.RoleCode).First(&role).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色已存在：" + toAddRole.RoleCode)
	}
	err = global.GqaDb.Create(&toAddRole).Error
	return err
}

func (s *ServiceRole) DeleteRoleById(id uint) (err error) {
	var sysRole model.SysRole
	if sysRole.Stable == "yesNo_yes" {
		return errors.New("系统内置不允许删除：" + sysRole.RoleCode)
	}
	if err = global.GqaDb.Where("id = ?", id).First(&sysRole).Error; err != nil {
		return err
	}
	roleCode := sysRole.RoleCode
	// 删除 sys_role_api 表的权限
	err = global.GqaDb.Where("role_code = ?", roleCode).Delete(model.SysRoleApi{}).Error
	// 删除 sys_role 表的数据
	err = global.GqaDb.Unscoped().Delete(&sysRole).Error
	if err != nil {
		return err
	}
	// 删除 sys_user_role 表的对应关系
	err = global.GqaDb.Where("sys_role_role_code = ?", roleCode).Delete(&model.SysUserRole{}).Error
	if err != nil {
		return err
	}
	// 删除 sys_role_menu 表的对应关系
	err = global.GqaDb.Where("sys_role_role_code = ?", roleCode).Delete(&model.SysRoleMenu{}).Error
	if err != nil {
		return err
	}
	return err
}

func (s *ServiceRole) QueryRoleById(id uint) (err error, roleInfo model.SysRole) {
	var role model.SysRole
	err = global.GqaDb.Preload("CreatedByUser").Preload("UpdatedByUser").
		Preload("DefaultPageMenu").First(&role, "id = ?", id).Error
	return err, role
}

func (s *ServiceRole) GetRoleMenuList(roleCode *model.RequestRoleCode) (err error, menu []model.SysRoleMenu) {
	err = global.GqaDb.Where("sys_role_role_code=?", roleCode.RoleCode).Find(&menu).Error
	return err, menu
}

func (s *ServiceRole) EditRoleMenu(toEditRoleMenu *model.RequestRoleMenuEdit) (err error) {
	err = global.GqaDb.Where("sys_role_role_code=?", toEditRoleMenu.RoleCode).Delete(&model.SysRoleMenu{}).Error
	if err != nil {
		return err
	}
	err = global.GqaDb.Where("sys_role_role_code=?", toEditRoleMenu.RoleCode).Delete(&model.SysRoleButton{}).Error
	if err != nil {
		return err
	}
	if len(toEditRoleMenu.RoleMenu) != 0 {
		err = global.GqaDb.Model(&model.SysRoleMenu{}).Create(&toEditRoleMenu.RoleMenu).Error
		if err != nil {
			return err
		}
	}
	if len(toEditRoleMenu.RoleButton) != 0 {
		err = global.GqaDb.Model(&model.SysRoleButton{}).Create(&toEditRoleMenu.RoleButton).Error
		if err != nil {
			return err
		}
	}
	defaultPage := toEditRoleMenu.DefaultPage
	err = global.GqaDb.Model(&model.SysRole{}).Where("role_code = ?", toEditRoleMenu.RoleCode).Update("default_page", defaultPage).Error
	return nil
}

func (s *ServiceRole) GetRoleApiList(roleCode *model.RequestRoleCode) (err error, api []model.SysRoleApi) {
	var roleApi []model.SysRoleApi
	err = global.GqaDb.Where("role_code=?", roleCode.RoleCode).Find(&roleApi).Error
	return err, roleApi
}

func (s *ServiceRole) EditRoleApi(toEditRoleApi *model.RequestEditRoleApi) (err error) {
	err = global.GqaDb.Where("role_code = ?", toEditRoleApi.RoleCode).Delete(&model.SysRoleApi{}).Error

	if len(toEditRoleApi.RoleApi) != 0 {
		err = global.GqaDb.Model(&model.SysRoleApi{}).Create(&toEditRoleApi.RoleApi).Error
		return err
	}
	return nil
}

func (s *ServiceRole) QueryUserByRole(roleCode *model.RequestRoleCode) (err error, user []model.SysUser) {
	role := model.SysRole{
		RoleCode: roleCode.RoleCode,
	}
	var userList []model.SysUser
	err = global.GqaDb.Model(&role).Association("User").Find(&userList)
	return err, userList
}

func (s *ServiceRole) RemoveRoleUser(toRemoveRoleUser *model.RequestRoleUser) (err error) {
	var roleUser model.SysUserRole
	if toRemoveRoleUser.Username == "admin" && toRemoveRoleUser.RoleCode == "super-admin" {
		return errors.New("抱歉，你不能把超级管理员从超级管理员组中移除！")
	}
	err = global.GqaDb.Where("sys_role_role_code = ? and sys_user_username = ?", toRemoveRoleUser.RoleCode, toRemoveRoleUser.Username).Delete(&roleUser).Error
	return err
}

func (s *ServiceRole) AddRoleUser(toAddRoleUser *model.RequestRoleUserAdd) (err error) {
	var roleUser []model.RequestRoleUser
	for _, r := range toAddRoleUser.Username {
		ur := model.RequestRoleUser{
			Username: r,
			RoleCode: toAddRoleUser.RoleCode,
		}
		roleUser = append(roleUser, ur)
	}
	if len(roleUser) != 0 {
		err = global.GqaDb.Model(&model.SysUserRole{}).Save(&roleUser).Error
		return err
	} else {
		return errors.New("本次操作没有影响！")
	}
}

func (s *ServiceRole) EditRoleDeptDataPermission(toEditRoleDeptDataPermission *model.RequestRoleDeptDataPermission) (err error) {
	var sysRole model.SysRole
	if err = global.GqaDb.Where("role_code = ?", toEditRoleDeptDataPermission.RoleCode).First(&sysRole).Error; err != nil {
		return err
	}
	if sysRole.Stable == "yesNo_yes" {
		return errors.New("系统内置不允许编辑：" + toEditRoleDeptDataPermission.RoleCode)
	}
	sysRole.DeptDataPermissionType = toEditRoleDeptDataPermission.DeptDataPermissionType
	sysRole.DeptDataPermissionCustom = toEditRoleDeptDataPermission.DeptDataPermissionCustom
	err = global.GqaDb.Save(&sysRole).Error
	return err
}

func (s *ServiceRole) GetRoleButtonList(roleCode *model.RequestRoleCode) (err error, button []model.SysRoleButton) {
	err = global.GqaDb.Where("sys_role_role_code=?", roleCode.RoleCode).Find(&button).Error
	return err, button
}
