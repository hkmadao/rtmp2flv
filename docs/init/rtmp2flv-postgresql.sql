-- public.client_info definition

-- Drop table

-- DROP TABLE public.client_info;

CREATE TABLE public.client_info (
	client_code varchar(255) NULL, -- 编号
	sign_secret varchar(255) NULL, -- 注册信息签名密钥
	secret varchar(255) NULL, -- 数据传输加密密钥
	note varchar(255) NULL, -- 备注
	id_client_info varchar(255) NOT NULL,
	CONSTRAINT pk_client_info PRIMARY KEY (id_client_info)
);
COMMENT ON TABLE public.client_info IS '客户端信息';

-- Column comments

COMMENT ON COLUMN public.client_info.client_code IS '编号';
COMMENT ON COLUMN public.client_info.sign_secret IS '注册信息签名密钥';
COMMENT ON COLUMN public.client_info.secret IS '数据传输加密密钥';
COMMENT ON COLUMN public.client_info.note IS '备注';


-- public.sys_token definition

-- Drop table

-- DROP TABLE public.sys_token;

CREATE TABLE public.sys_token (
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

-- DROP TABLE public.sys_user;

CREATE TABLE public.sys_user (
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


-- public.camera definition

-- Drop table

-- DROP TABLE public.camera;

CREATE TABLE public.camera (
	code varchar(255) NULL, -- 编号
	rtmp_auth_code varchar(255) NULL, -- rtmp识别码
	play_auth_code varchar(255) NULL, -- 播放权限码
	online_status bool NULL, -- 在线状态
	enabled bool NULL, -- 启用状态
	save_video bool NULL, -- 保存录像状态
	live bool NULL, -- 直播状态
	created timestamp NULL, -- 创建时间
	id varchar(255) NOT NULL,
	fg_encrypt bool NULL DEFAULT false, -- 加密标志
	fg_passive bool NULL DEFAULT false, -- 被动推送rtmp标志
	id_client_info varchar NULL,
	CONSTRAINT pk_camera PRIMARY KEY (id),
	CONSTRAINT camera_fk FOREIGN KEY (id_client_info) REFERENCES public.client_info(id_client_info)
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
COMMENT ON COLUMN public.camera.fg_encrypt IS '加密标志';
COMMENT ON COLUMN public.camera.fg_passive IS '被动推送rtmp标志';


-- public.camera_record definition

-- Drop table

-- DROP TABLE public.camera_record;

CREATE TABLE public.camera_record (
	id_camera_record varchar(255) NOT NULL,
	created timestamp NULL, -- 创建时间
	temp_file_name varchar(255) NULL, -- 临时文件名称
	fg_temp bool NULL, -- 临时文件标志
	file_name varchar(255) NULL, -- 文件名称
	fg_remove bool NULL, -- 文件删除标志
	duration int4 NULL, -- 文件时长: 单位：毫秒
	start_time timestamp NULL, -- 开始时间
	end_time timestamp NULL, -- 结束时间
	id_camera varchar(255) NOT NULL, -- 摄像头主属性
	has_audio bool NULL DEFAULT true, -- 是否有音频
	CONSTRAINT pk_camera_record PRIMARY KEY (id_camera_record),
	CONSTRAINT camera_record_fk FOREIGN KEY (id_camera) REFERENCES public.camera(id)
);
COMMENT ON TABLE public.camera_record IS '摄像头记录';

-- Column comments

COMMENT ON COLUMN public.camera_record.created IS '创建时间';
COMMENT ON COLUMN public.camera_record.temp_file_name IS '临时文件名称';
COMMENT ON COLUMN public.camera_record.fg_temp IS '临时文件标志';
COMMENT ON COLUMN public.camera_record.file_name IS '文件名称';
COMMENT ON COLUMN public.camera_record.fg_remove IS '文件删除标志';
COMMENT ON COLUMN public.camera_record.duration IS '文件时长: 单位：毫秒';
COMMENT ON COLUMN public.camera_record.start_time IS '开始时间';
COMMENT ON COLUMN public.camera_record.end_time IS '结束时间';
COMMENT ON COLUMN public.camera_record.id_camera IS '摄像头主属性';
COMMENT ON COLUMN public.camera_record.has_audio IS '是否有音频';


-- public.camera_share definition

-- Drop table

-- DROP TABLE public.camera_share;

CREATE TABLE public.camera_share (
	"name" varchar(255) NULL, -- 名称
	auth_code varchar(255) NULL, -- 权限码
	enabled bool NULL, -- 启用状态
	created timestamp NULL, -- 创建时间
	start_time timestamp NULL, -- 开始时间
	deadline timestamp NULL, -- 结束时间
	camera_id varchar(255) NOT NULL, -- 摄像头id
	id varchar(255) NOT NULL,
	CONSTRAINT pk_camera_share PRIMARY KEY (id),
	CONSTRAINT camera_share_camera_id_fkey FOREIGN KEY (camera_id) REFERENCES public.camera(id)
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