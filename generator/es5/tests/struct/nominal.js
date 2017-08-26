$module('nominal', function () {
  var $static = this;
  this.$type('099adf11', 'CoolBool', false, '', function () {
    var $instance = this.prototype;
    var $static = this;
    this.$box = function ($wrapped) {
      var instance = new this();
      instance[BOXED_DATA_PROPERTY] = $wrapped;
      return instance;
    };
    this.$roottype = function () {
      return $global.Boolean;
    };
    this.$typesig = function () {
      return {
      };
    };
  });

  this.$struct('0d40ca4a', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (someField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        someField: someField,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'someField', 'someField', function () {
      return $g.nominal.CoolBool;
    }, function () {
      return $g.________testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|2549c819<0d40ca4a>": true,
        "equals|4|2549c819<f361570c>": true,
        "Stringify|2|2549c819<ec87fc3f>": true,
        "Mapping|2|2549c819<95776681<any>>": true,
        "Clone|2|2549c819<0d40ca4a>": true,
        "String|2|2549c819<ec87fc3f>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
    var $result;
    var c;
    var s;
    var s2;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            c = $t.box($t.fastbox(true, $g.________testlib.basictypes.Boolean), $g.nominal.CoolBool);
            s = $g.nominal.SomeStruct.new(c);
            $promise.maybe($g.nominal.SomeStruct.Parse($g.________testlib.basictypes.JSON)($t.fastbox('{"someField": true}', $g.________testlib.basictypes.String))).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            s2 = $result;
            $resolve($t.fastbox(s2.someField.$wrapped && s.someField.$wrapped, $g.________testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
