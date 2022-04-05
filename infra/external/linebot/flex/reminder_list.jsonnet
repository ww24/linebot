local list = import 'reminder_list.json';

{
  type: 'carousel',
  contents: [
    {
      type: 'bubble',
      size: 'micro',
      header: {
        type: 'box',
        layout: 'vertical',
        contents: [
          {
            type: 'box',
            layout: 'horizontal',
            contents: [
              {
                type: 'box',
                layout: 'vertical',
                contents: [
                  {
                    type: 'text',
                    text: '✕',
                    weight: 'bold',
                    position: 'relative',
                    color: '#ffffff',
                    align: 'center',
                  },
                ],
                height: '24px',
                width: '24px',
                action: {
                  type: 'postback',
                  label: 'delete',
                  data: item.deleteTarget,
                },
              },
            ],
            position: 'relative',
            width: '100%',
            justifyContent: 'flex-end',
          },
          {
            type: 'text',
            text: item.title,
            color: '#ffffff',
            align: 'start',
            size: 'md',
            gravity: 'center',
          },
          {
            type: 'text',
            text: item.subTitle,
            color: '#ffffff',
            align: 'start',
            size: 'xs',
            gravity: 'center',
            margin: 'xs',
          },
        ],
        backgroundColor: '#27ACB2',
        paddingTop: '0px',
        paddingBottom: '10px',
        paddingStart: '10px',
        paddingEnd: '0px',
      },
      body: {
        type: 'box',
        layout: 'vertical',
        contents: [
          {
            type: 'box',
            layout: 'horizontal',
            contents: [
              {
                type: 'text',
                text: 'Next: ' + item.next,
                color: '#8C8C8C',
                size: 'sm',
                wrap: true,
              },
            ],
            flex: 1,
          },
        ],
        spacing: 'md',
        paddingAll: '12px',
      },
      styles: {
        header: {
          separator: false,
        },
        footer: {
          separator: false,
        },
      },
    }
    for item in list
  ],
}
