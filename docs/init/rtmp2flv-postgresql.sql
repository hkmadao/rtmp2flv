-- public.camera definition

-- Drop table

-- DROP TABLE camera;

CREATE TABLE camera (
	code varchar(255) NULL, -- 编号
	rtmp_auth_code varchar(255) NULL, -- rtmp识别码
	play_auth_code varchar(255) NULL, -- 播放权限码
	online_status bool NULL, -- 在线状态
	enabled bool NULL, -- 启用状态
	save_video bool NULL, -- 保存录像状态
	live bool NULL, -- 直播状态
	created timestamp NULL, -- 创建时间
	id varchar(255) NOT NULL,
	CONSTRAINT pk_camera PRIMARY KEY (id)
);
COMMENT ON TABLE public.camera IS '摄像头';

-- Column comments

COMMENT ON COLUMN public.camera.code IS '编号';
COMMENT ON COLUMN public.camera.rtmp_auth_code IS 'rtmp识别码';
COMMENT ON COLUMN public.camera.play_auth_code IS '播放权限码';
COMMENT ON COLUMN public.camera.online_status IS '在线状态';
COMMENT ON COLUMN public.camera.enabled IS '启用状态';
COMMENT ON COLUMN public.camera.save_video IS '保存录像状态';
COMMENT ON COLUMN public.camera.live IS '直播状态';
COMMENT ON COLUMN public.camera.created IS '创建时间';


-- public.sys_token definition

-- Drop table

-- DROP TABLE sys_token;

CREATE TABLE sys_token (
	username varchar(255) NULL, -- 用户名称
	nick_name varchar(255) NULL, -- 昵称
	create_time timestamp NULL, -- 创建时间
	"token" varchar(255) NULL, -- 令牌
	expired_time timestamp NULL, -- 过期时间
	user_info_string varchar(4000) NULL, -- 用户信息序列化
	id_sys_token varchar(255) NOT NULL,
	CONSTRAINT pk_sys_token PRIMARY KEY (id_sys_token)
);
COMMENT ON TABLE public.sys_token IS '令牌';

-- Column comments

COMMENT ON COLUMN public.sys_token.username IS '用户名称';
COMMENT ON COLUMN public.sys_token.nick_name IS '昵称';
COMMENT ON COLUMN public.sys_token.create_time IS '创建时间';
COMMENT ON COLUMN public.sys_token."token" IS '令牌';
COMMENT ON COLUMN public.sys_token.expired_time IS '过期时间';
COMMENT ON COLUMN public.sys_token.user_info_string IS '用户信息序列化';


-- public.sys_user definition

-- Drop table

-- DROP TABLE sys_user;

CREATE TABLE sys_user (
	account varchar(255) NULL, -- 登录账号 
	user_pwd varchar(255) NULL, -- 用户密码 
	phone varchar(255) NULL, -- 手机号码
	email varchar(255) NULL, -- 邮箱
	"name" varchar(255) NULL, -- 姓名 
	nick_name varchar(255) NULL, -- 昵称
	gender varchar(255) NULL, -- 性别
	fg_active bool NULL, -- 启用标志
	id_user varchar(255) NOT NULL,
	CONSTRAINT pk_sys_user PRIMARY KEY (id_user)
);
COMMENT ON TABLE public.sys_user IS '系统用户';

-- Column comments

COMMENT ON COLUMN public.sys_user.account IS '登录账号 ';
COMMENT ON COLUMN public.sys_user.user_pwd IS '用户密码 ';
COMMENT ON COLUMN public.sys_user.phone IS '手机号码';
COMMENT ON COLUMN public.sys_user.email IS '邮箱';
COMMENT ON COLUMN public.sys_user."name" IS '姓名 ';
COMMENT ON COLUMN public.sys_user.nick_name IS '昵称';
COMMENT ON COLUMN public.sys_user.gender IS '性别';
COMMENT ON COLUMN public.sys_user.fg_active IS '启用标志';


-- public.camera_share definition

-- Drop table

-- DROP TABLE camera_share;

CREATE TABLE camera_share (
	"name" varchar(255) NULL, -- 名称
	auth_code varchar(255) NULL, -- 权限码
	enabled bool NULL, -- 启用状态
	created timestamp NULL, -- 创建时间
	start_time timestamp NULL, -- 开始时间
	deadline timestamp NULL, -- 结束时间
	camera_id varchar(255) NOT NULL, -- 摄像头id
	id varchar(255) NOT NULL,
	CONSTRAINT pk_camera_share PRIMARY KEY (id),
	CONSTRAINT camera_share_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES camera(id)
);
COMMENT ON TABLE public.camera_share IS '摄像头分享';

-- Column comments

COMMENT ON COLUMN public.camera_share."name" IS '名称';
COMMENT ON COLUMN public.camera_share.auth_code IS '权限码';
COMMENT ON COLUMN public.camera_share.enabled IS '启用状态';
COMMENT ON COLUMN public.camera_share.created IS '创建时间';
COMMENT ON COLUMN public.camera_share.start_time IS '开始时间';
COMMENT ON COLUMN public.camera_share.deadline IS '结束时间';
COMMENT ON COLUMN public.camera_share.camera_id IS '摄像头id';