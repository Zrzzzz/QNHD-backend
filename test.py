l = """floor_report_solve
post_valued
been_blocked
post_department_transfer
post_report_solve
post_deleted_by_report
floor_deleted_by_report
post_type_transfer"""
a = """您好，您在<post>下举报的评论<floor>已被移除。感谢您对论坛秩序的维护。
恭喜微友，您的<post>被加为精华帖！
您好，您因<reason>被禁言<day>天。请您遵守社区规范发言，感谢您的配合。
您在校务专区发布的帖子<post>，不属于所选部门问题，已被移交到<department>部门下。感谢您的支持。
您好，您举报的帖子<post>已被移除。感谢您对论坛秩序的维护。
您好，您在“<post>”下的评论“<floor>”因“<reason>”已被移除。请您遵守社区规范发言，感谢您的配合。
您在“<from_type>”发布的帖子“<post>”，不属于该分区相关问题，已被转移至<to_type>由大家一起讨论。感谢您的支持。
您好，您的帖子<post>因<reason>已被移除。请您遵守社区规范发帖，感谢您的配合。"""
import re
for i, type in enumerate(l.split('\n')):
    st = a.split('\n')[i]
    ret = re.findall(r'<(.+?)>', st)
    print('NoticeType.{}: {}, '.format(type.upper(), ret).replace('\'', '"').replace("[", "{").replace("]", "}"))