$module('nullbox', function () {
  var $static = this;
  this.$struct('45b3c939', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'Value', 'Value', function () {
      return $g.________testlib.basictypes.String;
    }, function () {
      return $g.________testlib.basictypes.String;
    }, true);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|fd8bc7c9<45b3c939>": true,
        "equals|4|fd8bc7c9<54ff3ddf>": true,
        "Stringify|2|fd8bc7c9<44e219a9>": true,
        "Mapping|2|fd8bc7c9<ad6de9ce<any>>": true,
        "Clone|2|fd8bc7c9<45b3c939>": true,
        "String|2|fd8bc7c9<44e219a9>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $temp0;
    var s;
    s = ($temp0 = $g.nullbox.SomeStruct.new(), $temp0.Value = null, $temp0);
    return $t.fastbox((s.Value == null) && (s.Mapping().$index($t.fastbox('Value', $g.________testlib.basictypes.String)) == null), $g.________testlib.basictypes.Boolean);
  };
});
